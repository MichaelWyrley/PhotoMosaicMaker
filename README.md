# PhotoMoseicMaker
This application takes an image, a search and a number of smaller images. It will then use them to create a photo mosaic of the image using google images.

# How the application works

1. Use a Google images scraper in order to download the requested images
2. Get the average colour value for each downloaded image
3. Go through the given image and replace each pixle (or group of pixles) with the closest matching downloaded image (using colour value)
4. Output the resulting image


## Image Scraper

I mostly referenced this blog post for getting images from google chrome [Blog Post](https://towardsdatascience.com/image-scraping-with-python-a96feda8af2d)

The image scraper takes in the query, number of images wanted and the location they are going to be stored.
I am using an anoconda environment so you can use
`conda activate photoMoseic` 
in order to install all dependencys

When you want to use the script please use this form
`python ./ImageScraper.py <QUERY> <NO IMAGES> <PATH>`

e.g.
`python ./ImageScraper.py cat 5 ./images/`

### Trouble Shooting

Check that your chrome version and the chrome driver version match, if they don't then download the one that matches your chrome version from [here](https://chromedriver.chromium.org/downloads)




