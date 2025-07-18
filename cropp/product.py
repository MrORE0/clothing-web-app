from pydantic import BaseModel, HttpUrl
from typing import List, Dict


class ProductVariant(BaseModel):
    color: str
    # images: List[HttpUrl]  # or just List[str] if URLs aren't validated
    sizes: Dict[str, bool]  # e.g. {'XS': False, 'S': True}
    url: HttpUrl
    price: float
    discounted_price: float | None
    currency: str


class Product(BaseModel):
    name: str
    variants: List[ProductVariant]
