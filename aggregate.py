import asyncio
from fast_ext import fast_ext
from ookla_ext import ookla_ext
from test_ext import testnet_ext 


async def aggregate(sites=None):

    # scrapers={
    #     'fast':fast_ext(),
    #     'ookla':ookla_ext(),
    #     'testnet':testnet_ext()
    # }
    # loop through dict to avoid if hell
    output = {}
    if sites:
        if 'fast' in sites:
            output['fast'] = await fast_ext()
        if 'ookla' in sites:
            output['ookla'] = await ookla_ext()
        if 'testnet' in sites:
            output['testnet'] = await testnet_ext()
      
    else:
        fast_task = fast_ext()
        ookla_task = ookla_ext()
        testnet_task = testnet_ext()

        fast, ookla, testnet = await asyncio.gather(fast_task, ookla_task, testnet_task)

        output = {
            'fast': fast,
            'ookla': ookla,
            'testnet': testnet
        }

    # print("Aggregate result ready.")
    return output

if __name__ == "__main__":
    result = asyncio.run(aggregate())
    print(result)
