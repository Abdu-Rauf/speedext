from playwright.async_api import async_playwright
import asyncio

async def fast_ext():
    async with async_playwright() as f:

        #Launch the browser
        browser = await f.firefox.launch()
        page = await browser.new_page()
        await page.goto('https://fast.com/')

        # Wait for speed test to complete
        await page.wait_for_selector('.speed-units-container.succeeded')
        # await asyncio.sleep(10)
        speed = await page.text_content('#speed-value')
        speed_unit = await page.text_content('#speed-units')

        # Close the Browser and return the result
        await browser.close()

        return speed + " " + speed_unit
    
async def main():
    result = await fast_ext()
    print(result)


if __name__ == '__main__':
    asyncio.run(main())