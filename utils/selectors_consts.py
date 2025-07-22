NAME_SELECTOR = "h1[data-testid*='product-name']"
SIZE_PICKER_SELECTOR = "ul[data-testid*='product-size-group'] li"
COLORS = "li[data-testid='color-picker-list-item'] button[data-testid='color-picker-list-item-button']"
PRICE = "div[data-selen='product-price']"
DISCOUNT = "div[data-selen='product-discount-price']"
REGULAR_PRICE = "div[data-selen='product-regular-price']"
COOKIE_BUTTON_ID = "#cookiebotDialogOkButton"
COOKIE_BANNER = "#cookiebanner"
COLOR_NAME_SELECTORS = {
    "reserved": "div[data-testid*='color-picker-color-name']",
    "default": "span[data-testid*='color-picker-color-name']",
}


def get_color_name_selector(url: str) -> str:
    if "reserved" in url:
        return COLOR_NAME_SELECTORS["reserved"]
    return COLOR_NAME_SELECTORS["default"]
