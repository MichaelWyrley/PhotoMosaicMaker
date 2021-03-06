# PhotoMoseicMaker
This application takes an image, a search and a number of smaller images. It will then use them to create a photo mosaic of the image using google images.

# How the application works

1. Use a Google images scraper in order to download the requested images
2. Get the average colour value for each downloaded image
3. Go through the given image and replace each pixle (or group of pixles) with the closest matching downloaded image (using colour value)
4. Output the resulting image


## Image Scraper

I mostly referenced this blog post for getting images from google chrome [Blog Post](https://towardsdatascience.com/image-scraping-with-python-a96feda8af2d)

I have created both a Google Image Scraper and a Free Image scraper which uses the pexels image library (each image has a free license)
The image scraper takes in the query, number of images wanted, the location they are going to be stored and the final dimensions.
I am using an anoconda environment so you can use
`conda activate photoMoseic` 
in order to install all dependencys

When you want to use the script please use this form
`python ./GoogleImageScraper.py <QUERY> <NO IMAGES> <PATH> <WIDTHxHEIGHT>`
`python ./FreeImageScraper.py <QUERY> <NO IMAGES> <PATH> <WIDTHxHEIGHT>`

e.g.
`python ./GoogleImageScraper.py 'cat' 5 './images/' '100x100'`
`python ./FreeImageScraper.py 'cat' 5 './images/' '100x100'`

## Mosaic generateor

The mosaic generator takes the image to be converted, the number of sample images you want to use, the scale of the sample images, the amount you want to shrink the image given and the final location of the image.

The image converter only works on .jpg or .jpeg images

It uses go with a specific library installed for resizing the image ([library](https://github.com/nfnt/resize))
use `go get github.com/nfnt/resize` if you don't have it installed

Run the program via the .exe
`./photoMosaic -img <IMAGE> -no <NO SAMPLE IMAGES>  -shrink <AMOUNT YOU WANT TO DECREASE BY> -location <LOCATION OF FINAL IMAGE> -images <LOCATION OF SAMPLE IMAGES>`

e.g.
`./photoMosaic -img '../cat1.jpeg' -no '10' -shrink 10 -location "../image" -images "../images/"`


## Trouble Shooting

Check that your chrome version and the chrome driver version match, if they don't then download the one that matches your chrome version from [here](https://chromedriver.chromium.org/downloads)

You might need to make the directory that the images to be scraped will be stored in before linking to it

# Example uses

Before (image from [Pixabay](https://www.pexels.com/@pixabay/))
<img src="./frog.jpg"></img>

After
<img src="./convertedFrog.jpg"></img>

# Disclaimer

If you use the Google Image Scraper. These images may or may not have copyright associated with them. Do not use these images in a way that would infringe on the original image holders copyright
