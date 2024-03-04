from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service as ChromeService
from fake_useragent import UserAgent
import time
import random
import csv

BASE_URL = 'https://www.tandem-associates.com/lionel/'

def setup_browser():
    """Set up the browser with desired options and return the browser instance."""
    ua = UserAgent()
    chrome_options = Options()
    chrome_options.add_argument(f"user-agent={ua.random}")  # Rotate user agent
    chrome_options.add_argument("--headless")  # Run in headless mode
    chrome_options.add_argument("--disable-web-security")
    chrome_options.add_argument("--allow-running-insecure-content")
    chrome_options.add_argument("--no-sandbox")
    
    return webdriver.Chrome(service=ChromeService(executable_path='../../../Downloads/chromedriver-mac-x64/chromedriver'), options=chrome_options)

def extract_image_tags(browser, link):
    """Extract <img> tags for a given link."""
    full_url = BASE_URL + link
    browser.get(full_url)
    time.sleep(random.uniform(1, 3))  # Random sleep time between 1 to 3 seconds to mimic human behavior
    soup = BeautifulSoup(browser.page_source, 'lxml')
    img_tags = soup.select('td[align="left"] img')
    return [str(tag) for tag in img_tags]

def main():
    browser = setup_browser()
    img_tags_list = []

    with open('train_links.csv', 'r') as csvfile:
        reader = csv.reader(csvfile)
        next(reader)  # Skip header row
        links = [row[1] for row in reader]

    for link in links:
        try:
            img_tags = extract_image_tags(browser, link)
            img_tags_list.extend(img_tags)
        except Exception as e:
            print(f"Error extracting data for {link}. Error: {e}")

    with open('train_images.csv', 'w', newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(["Image Tag"])
        for img_tag in img_tags_list:
            writer.writerow([img_tag])

    browser.quit()

if __name__ == "__main__":
    main()
