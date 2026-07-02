import asyncio
import io
import json
import logging
import re
import pytesseract
from bs4 import BeautifulSoup
from PIL import Image
from playwright.async_api import async_playwright

#the path will be changed in Docker
pytesseract.pytesseract.tesseract_cmd = '/opt/homebrew/bin/tesseract'

# --- CONFIGURATION ---
LISTING_URL = "https://www.topjobs.lk/applicant/vacancybyfunctionalarea.jsp?FA=&jst=OPEN&sQut=&txtKeyWord=&chkGovt=&chkParttime=&chkWalkin=&chkNGO="
TESSERACT_LANG = "eng+sin+tam"
POPUP_WAIT_TIMEOUT = 15000 
USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
log = logging.getLogger("topjobs_crawler")

def get_total_pages(html: str) -> int:
    soup = BeautifulSoup(html, "html.parser")
    pagination_div = soup.find("div", class_="pagin-block page-show")
    
    if pagination_div:
        text = pagination_div.get_text(strip=True)
        # Regex to find the number before "page(s)"
        match = re.search(r'(\d+)\s*page\(s\)', text)
        if match:
            return int(match.group(1))
    
    return 1 # Default to 1 

def parse_listing_html(html: str) -> list[dict]:
    soup = BeautifulSoup(html, "html.parser")
    jobs = []
    rows = soup.find_all("tr", attrs={"onclick": re.compile(r"createAlert")})
    
    for row in rows:
        row_id = row.get("id")
        tds = row.find_all("td", recursive=False)
        if len(tds) < 6: continue
        
        jobs.append({
            "row_id": row_id,
            "title": tds[2].find("h2").text.strip() if tds[2].find("h2") else "N/A",
            "employer": tds[2].find("h1").text.strip() if tds[2].find("h1") else "N/A",
            "location": tds[6].text.strip() if len(tds) > 6 else (tds[5].text.strip() if len(tds) > 5 else "N/A"),
            "ocr_text": None,
            "error": None
        })
    return jobs

async def run_crawler():
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)
        context = await browser.new_context(user_agent=USER_AGENT)
        page = await context.new_page()

        await page.goto(f"{LISTING_URL}&pageNo=1", wait_until="networkidle")
        total_pages = get_total_pages(await page.content())
        log.info(f"Total pages detected: {total_pages}")


        all_jobs = []
        
        for page_num in range(1, total_pages + 1):
            # log.info("Navigating to listings...")
            # await page.goto(LISTING_URL, wait_until="networkidle")

            log.info(f"Scraping page {page_num} of {total_pages}...")
            if page_num > 1:
                await page.goto(f"{LISTING_URL}&pageNo={page_num}", wait_until="networkidle")
            
            content = await page.content()
            jobs = parse_listing_html(content)
            log.info(f"Found {len(jobs)} jobs. Starting popup processing...")

            x = 0

            for i, job in enumerate(jobs):
                try:
                    log.info(f"[{i+1}/{len(jobs)}] Processing {job['employer']}")

                    await asyncio.sleep(3)
                    
                    async with context.expect_page(timeout=POPUP_WAIT_TIMEOUT) as popup_info:
                        await page.evaluate(f"document.getElementById('{job['row_id']}').click()")
                    
                    popup = await popup_info.value
                    await popup.wait_for_load_state("networkidle")

                    # --- NEW LOGIC: Target the large advertisement image ---
                    # We find all images and pick the one with a width > 500px
                    all_images = popup.locator("img")
                    img_locator = None
                    
                    for count in range(await all_images.count()):
                        candidate = all_images.nth(count)
                        box = await candidate.bounding_box()
                        if box and box['width'] > 500:
                            img_locator = candidate
                            break
                    
                    if not img_locator:
                        log.error("Could not find a large advertisement image.")
                        await popup.close()
                        continue
                        #raise Exception("Could not find a large advertisement image.")

                    await img_locator.wait_for(state="visible", timeout=POPUP_WAIT_TIMEOUT)
                    screenshot_bytes = await img_locator.screenshot()
                    await popup.close()

                    # --- OCR Processing ---
                    image = Image.open(io.BytesIO(screenshot_bytes)).convert("L")
                    if image.width < 1400:
                        scale = 1400 / image.width
                        image = image.resize((int(image.width * scale), int(image.height * scale)))
                    
                    image = image.point(lambda x: 0 if x < 180 else 255, '1')
                    #job["ocr_text"] = pytesseract.image_to_string(image, lang=TESSERACT_LANG).strip()
                    job["ocr_text"] = " ".join(pytesseract.image_to_string(image, lang=TESSERACT_LANG).split())

                    all_jobs.append(job)

                    x = x + 1
                    if x >= 5:
                        await asyncio.sleep(10)
                        break
                    
                except Exception as e:
                    log.error(f"Failed to process {job['row_id']}: {e}")
                    job["error"] = str(e)

        await browser.close()
        
        # with open("vacancies.json", "w", encoding="utf-8") as f:
        #     json.dump(jobs, f, ensure_ascii=False, indent=2)
        print(all_jobs)
        log.info("Done! Data saved to vacancies.json")
        return all_jobs
