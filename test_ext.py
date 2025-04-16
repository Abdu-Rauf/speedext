import asyncio
from playwright.async_api import async_playwright

async def testnet_ext():
    async with async_playwright() as f:

        # Launch the Browser
        browser = await f.firefox.launch()
        page = await browser.new_page()

        #Navigate to testynet and click on test my download
        await page.goto('https://testmy.net/',wait_until='domcontentloaded')
        await page.click('[href="/download"]')

        #Wait for Run test button to apppear
        await asyncio.sleep(2)

        await page.click('.button.btn.btnGr.btn-default.small.button22')
        await asyncio.sleep(15)

        #Extract speed value and close
        await page.wait_for_selector('.color22')
        speed = await page.text_content('.color22')
        await browser.close()

        return {'testmynet':speed}


async def main():
    result = await testnet_ext()
    print(result)


if __name__ == '__main__':
    asyncio.run(main())