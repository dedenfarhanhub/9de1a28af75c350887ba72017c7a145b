'''
this module runs golint on all Go scripts found in a directory tree
'''
import asyncio
import os
import platform
import sys
from concurrent.futures.thread import ThreadPoolExecutor

mag_count = 0

async def check_async(output_file, module, semaphore):
    '''apply golint to the file specified if it is a *.go file'''
    global mag_count
    if module.endswith(".go") and not (module.endswith("bigquery.go") or module.endswith("lint.py")):
        async with semaphore:
            print(f"CHECKING {module}")
            command = f"golint {module}"
            process = await asyncio.create_subprocess_shell(command, stdout=asyncio.subprocess.PIPE,
                                                            stderr=asyncio.subprocess.PIPE)
            stdout, stderr = await process.communicate()
            data = stdout.decode('utf-8')
            for line in data.splitlines():
                if line and "don't use ALL_CAPS" not in line:
                    mag_count += 1
            return data

async def process_files_async(base_directory, output_file, semaphore):
    coroutines = []
    for root, dirs, files in os.walk(base_directory):
        for name in files:
            filepath = os.path.join(root, name)
            if 'vendor' not in filepath:
                coroutines.append(check_async(output_file, filepath, semaphore))

    results = await asyncio.gather(*coroutines)  # Use asyncio.gather() instead of asyncio.wait()

    with open(output_file, 'a', encoding='utf-8') as infile:
        for result in results:
            if result:
                infile.write(result)
                print(result)

if __name__ == "__main__":
    BASE_DIRECTORY = os.getcwd()
    try:
        print(sys.argv)
        OUTPUT_FILE = sys.argv[1]
    except IndexError:
        OUTPUT_FILE = 'golint.log'

    print("looking for *.go scripts in subdirectories of ", BASE_DIRECTORY)

    # Set up the event loop based on the operating system
    if platform.system() == 'Windows':
        loop = asyncio.ProactorEventLoop()
        asyncio.set_event_loop(loop)
    else:
        loop = asyncio.get_event_loop()

    executor = ThreadPoolExecutor(max_workers=2)
    loop.set_default_executor(executor)
    semaphore = asyncio.Semaphore(4)  # Removed loop=loop here

    loop.run_until_complete(process_files_async(BASE_DIRECTORY, OUTPUT_FILE, semaphore))
    loop.close()

    print("Done linting!")
    print("==" * 50)
    print(f"{mag_count} errors found")
