from pydantic import BaseModel, field_validator


class RawJobInput(BaseModel):
    employer: str = ""
    job_role: str
    location: str = ""
    description: str = ""
    crawler_run_id: int
    source: str

    @field_validator("crawler_run_id")
    @classmethod
    def crawler_run_id_must_be_positive(cls, v: int) -> int:
        if v <= 0:
            raise ValueError("crawler_run_id must be greater than 0")
        return v
 
    @field_validator("source")
    @classmethod
    def source_must_not_be_blank(cls, v: str) -> str:
        if not v or not v.strip():
            raise ValueError("source must not be blank")
        return v