from pydantic import BaseModel, Field


class RawJobInput(BaseModel):
    employer: str = Field(default= "")
    job_role: str = Field(default= "")
    location: str = Field(default= "")
    description: str = Field(default= "")
    crawler_run_id: int = Field(gt=0)
    source: str = Field(pattern=r"\S")