import requests
import pytesseract
from PIL import Image
from io import BytesIO
from bs4 import BeautifulSoup
from crawl4ai.async_webcrawler import AsyncWebCrawler

# Configure Tesseract path if it's not in your system PATH
# pytesseract.pytesseract.tesseract_cmd = r'/usr/bin/tesseract'

TESSERACT_LANG = "eng+sin+tam"

def remove_sinhala_control_chars(text):
    cleaned_text = text.replace('\u200c', '').replace('\u200d', '')
    return cleaned_text

def perform_ocr(image_url):
    """Downloads image and extracts text."""
    try:
        response = requests.get(image_url, timeout=10)
        img = Image.open(BytesIO(response.content)).convert('L') # Convert to grayscale for better OCR
        return " ".join(pytesseract.image_to_string(img, lang=TESSERACT_LANG).split())
    except Exception as e:
        return f"OCR Error: {e}"

async def fetch_job_details():
    base_url = "https://governmentjobs.lk/index.php?page={}&ipp=25&"
    jobs_data = []
    
    # Open the crawler once for the whole process
    async with AsyncWebCrawler() as crawler:
        # 1. Fetch the main page
        result = await crawler.arun(url="https://governmentjobs.lk/index.php?page=1&ipp=25&")
        soup = BeautifulSoup(result.html, 'html.parser')

        pagination_text = soup.select_one('.category-results .paginate').text
        total_pages = int(pagination_text.split()[-1])
        print(total_pages)
        
        for page_num in range(1, total_pages + 1):
            current_url = base_url.format(page_num)
            print(f"Crawling: {current_url}")
            
            # Fetch the specific page
            page_result = await crawler.arun(url=current_url)
            page_soup = BeautifulSoup(page_result.html, 'html.parser')
            # 2. Iterate through each job
            for job_div in page_soup.select('.grid-view.product'):
                try:
                    title = job_div.select_one('h5 strong').text.strip()
                    employer = job_div.select_one('a[href*="vtag"]').text.strip()
                    
                    # Get the link to the detailed image page
                    image_page_link = job_div.select_one('a[href*="image-view.php"]')['href']
                    full_image_page_url = "https://governmentjobs.lk/" + image_page_link
                    
                    # 3. Visit the detail page using the same crawler session
                    # MUST use await and arun()
                    img_res = await crawler.arun(url=full_image_page_url)
                    img_soup = BeautifulSoup(img_res.html, 'html.parser')
                    # img_tag = img_soup.select_one('.page-content img')
                    
                    # description = ""
                    # if img_tag and 'src' in img_tag.attrs:
                    #     description = perform_ocr(img_tag['src'])

                    img_tags = img_soup.select('.page-content img')

                    description = ""
                    for img_tag in img_tags:
                        src = img_tag.get('src')
                        if src and "amazonaws.com/mytutor.lk/vacancy" in src:
                            description = description + perform_ocr(src)

                    description = remove_sinhala_control_chars(description)
                    
                    jobs_data.append({
                        "title": title,
                        "employer": employer,
                        "location": "Sri Lanka",
                        "description": description.strip()
                    })
                    
                except Exception as e:
                    print(f"Error parsing job: {e}")
                    continue
            
    return jobs_data