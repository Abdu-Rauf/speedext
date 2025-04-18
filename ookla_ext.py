# speed_class ='result-data-large number result-data-value download-speed'

# button_class = 'start-text'

# unit = 'result-data-unit'
import asyncio
from playwright.async_api import async_playwright

async def ookla_ext():
    async with async_playwright() as f:

        #Launch the browser
        browser = await f.chromium.launch(
            # headless=False
        )
        # context = await browser.new_context(
        #     user_agent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        #     viewport={"width": 1280, "height": 720},
        # )

        page = await browser.new_page()

        #         await page.add_init_script("""
        # Object.defineProperty(navigator, 'webdriver', { get: () => false });
        # """)

        await page.goto('https://www.speedtest.net/')



        # Click GO
        await page.wait_for_selector("span.start-text", timeout=10000)
        # await page.hover("span.start-text")
        await page.click('span.start-text')
        await asyncio.sleep(30) # Wait for speed test to finitsh completing

        
        speed_div = await page.wait_for_selector('span.result-data-large.number.result-data-value.download-speed', timeout = 35000)
        unit_div = await page.wait_for_selector('span.result-data-unit')

        speed = await speed_div.text_content()
        speed_unit = await unit_div.text_content()

        
        return speed+ " " +speed_unit
    
async def main():
    result = await ookla_ext()
    print(result)


if __name__ == '__main__':
    asyncio.run(main())

        