import asyncio
from fast_ext import fast_ext
from ookla_ext import ookla_ext


async def aggregate(sites=None):

    output = {}
    if sites:
        if 'fast' in sites:
            output['fast'] = await fast_ext()
        if 'ookla' in sites:
            output['ookla'] = await ookla_ext()
      
    else:
        fast_task = fast_ext()
        ookla_task = ookla_ext()

        fast, ookla, testnet = await asyncio.gather(fast_task, ookla_task)

        output = {
            'fast': fast,
            'ookla': ookla,
        }

    return output

if __name__ == "__main__":
    result = asyncio.run(aggregate())
    print(result)
