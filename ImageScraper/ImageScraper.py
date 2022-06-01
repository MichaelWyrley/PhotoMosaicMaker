from selenium import webdriver
from selenium.webdriver.common.by import By
from PIL import Image
import sys
import time
import requests
import io

# Where the Chrome Driver is stored
DRIVER_PATH = './chromedriver.exe'
DELAY = 2
PATH = "../images/"
NO_IMAGES = 20
IMAGE_TEXT = "dog"
IMAGE_DIMENSIONS = (400,400)


# Scroll to the bottom of a webpage
def scroll_to_end(wd, delay):
    wd.execute_script("window.scrollTo(0, document.body.scrollHeight);")
    time.sleep(delay)

# Get all the image urls so that they can be downloaded
def fetch_image_urls(query, number_of_images, wd, delay):

    # Make the Google Query
    search_url = "https://www.google.com/search?safe=off&site=&tbm=isch&source=hp&q={q}&oq={q}&gs_l=img".format(q=query)
    wd.get(search_url)

    image_urls = []
    image_count = 0
    results_start = 0
    previous_results = 0

    while image_count < number_of_images:
        scroll_to_end(wd, delay)

        # get all image thumbnail results
        thumbnail_results = wd.find_elements(By.CLASS_NAME, value="Q4LuWd")
        number_results = len(thumbnail_results)
        
        print(f"Found: {number_results} search results. Extracting links from {results_start}:{number_results}")
        
        # Click load more images if you run out of images
        if previous_results == number_results:
            print("Need to load more images")
            time.sleep(delay)
            wd.find_element(By.CLASS_NAME, value="mye4qd").click()
        
        #Loop through all the results that have been rechrieved (this is a slice so we don't repeat results when scrolling down a page)
        for img in thumbnail_results[results_start : number_results]:
            # Extract image thumbnail URLs
            attribute = img.get_attribute('src')
            if attribute and 'http' in attribute:
                image_urls.append(attribute)
                image_count += 1

            if image_count >= number_of_images:
                print("Found all images")
                break 
    
        results_start = len(thumbnail_results)    
        previous_results = number_results  
    
    return image_urls
                    
def download_image(path, urls, image_dimensions):

    for i, url in enumerate(urls):
        try:
            # Open the image
            image_content = requests.get(url).content
            image_file = io.BytesIO(image_content)
            image = Image.open(image_file)
            image = image.resize(image_dimensions)

            # Save the image
            file_path = path  + str(i) + ".jpg"
            with open(file_path, "wb") as f:
                image.save(f, "JPEG")

            print("Saved image ", file_path)

        except Exception as e:
            print(f"ERROR - Could not save {url} - {e}")




def scrape_images(query=IMAGE_TEXT, no_images=NO_IMAGES, download_path=PATH, image_dimensions=IMAGE_DIMENSIONS):
    wd = webdriver.Chrome(executable_path=DRIVER_PATH)
    wd.get('https://google.com')


    # Get rid of the popup window when using chrome
    wd.find_element(By.ID, value='L2AGLb').click()

    urls = fetch_image_urls(query, no_images, wd, DELAY)
    download_image(download_path, urls, image_dimensions)

    wd.quit()

args = sys.argv[1:]
query = args[0]
no_images = int(args[1])
path = args[2]
image_dimensions = tuple(map(int, args[3].split('x')))

scrape_images(query, no_images, path, image_dimensions)