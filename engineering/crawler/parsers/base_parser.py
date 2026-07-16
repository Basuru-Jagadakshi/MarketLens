from abc import ABC, abstractmethod

class BaseJobParser(ABC):

    @abstractmethod
    def parse_rule_based_fields(self, markdown: str) -> dict:
        pass