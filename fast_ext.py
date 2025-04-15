from playwright.async_api import async_playwright
import asyncio

async def fast_ext():
    async with async_playwright() as f:

        browser = await f.firefox.launch()
        page = await browser.new_page()
        await page.goto('https://fast.com/')
        await page.wait_for_selector('.speed-units-container.succeeded')
        await asyncio.sleep(10)
        speed = await page.text_content('#speed-value')
        speed_unit = await page.text_content('#speed-units')

        await browser.close()

        return { 'speed': speed + " " + speed_unit}
    
async def main():
    result = await fast_ext()
    print(result)


if __name__ == '__main__':
    asyncio.run(main())