from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service as ChromeService
from fake_useragent import UserAgent
import csv
import time
import random

def setup_browser():
    """Setup browser with desired options and return the browser instance."""
    ua = UserAgent()
    chrome_options = Options()
    chrome_options.add_argument(f"user-agent={ua.random}")

    return webdriver.Chrome(service=ChromeService(executable_path='../../../Downloads/chromedriver-mac-x64/chromedriver'), options=chrome_options)


def extract_links_from_column(column):
    """Extract links from a given column."""
    links = column.select('a')
    data = []

    for link in links:
        href = link['href']
        text = link.get_text().strip()

        # Check for additional string
        next_sibling = link.find_next_sibling(text=True)
        if next_sibling and not next_sibling.isspace():
            text += '-' + next_sibling.strip()

        data.append([text, href])

    return data


def scrape_train_links(browser):
    """Scrape train links from the page."""
    time.sleep(random.randint(2, 8))
    soup = BeautifulSoup(browser.page_source, 'lxml')

    columns = soup.select('td')
    all_data = []

    for column in columns:
        column_data = extract_links_from_column(column)
        all_data.extend(column_data)

    return all_data


def export_to_csv(data):
    """Export the scraped data to CSV."""
    with open('train_links.csv', 'w', newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(['Title', 'Link'])  # Header row
        writer.writerows(data)


def main():
    browser = setup_browser()
    try:
        browser.get('https://www.tandem-associates.com/lionel/lionel_trains_master_index.htm#TOPLIST')
        data = scrape_train_links(browser)
        export_to_csv(data)
    except Exception as e:
        print(f"An error occurred: {e}")
    finally:
        browser.quit()


if __name__ == "__main__":
    main()


# from bs4 import BeautifulSoup
# from selenium import webdriver
# from selenium.webdriver.chrome.options import Options
# from selenium.webdriver.chrome.service import Service as ChromeService
# from fake_useragent import UserAgent
# import csv
# import time
# import random
# import pandas as pd

# def setup_browser():
#     """Setup browser with desired options and return the browser instance."""
#     ua = UserAgent()
#     chrome_options = Options()
#     chrome_options.add_argument(f"user-agent={ua.random}")

#     return webdriver.Chrome(service=ChromeService(executable_path='../../../Downloads/chromedriver-mac-x64/chromedriver'), options=chrome_options)


# def scrape_train_links(browser):
#     """Scrape train links from the page."""
#     time.sleep(random.randint(2, 8))
#     soup = BeautifulSoup(browser.page_source, 'lxml')
    
#     links = soup.select('td[width="111"] a')
#     data = []

#     for link in links:
#         href = link['href']
#         text = link.get_text().strip()

#         # Check for additional string
#         next_sibling = link.find_next_sibling(text=True)
#         if next_sibling and not next_sibling.isspace():
#             text += '-' + next_sibling.strip()

#         data.append([text, href])

#     return data


# def export_to_csv(data):
#     """Export the scraped data to CSV."""
#     with open('train_links.csv', 'w', newline='') as csvfile:
#         writer = csv.writer(csvfile)
#         writer.writerow(['Title', 'Link'])  # Header row
#         writer.writerows(data)


# def main():
#     browser = setup_browser()
#     try:
#         browser.get('https://www.tandem-associates.com/lionel/lionel_trains_master_index.htm#TOPLIST')
#         data = scrape_train_links(browser)
#         export_to_csv(data)
#     except Exception as e:
#         print(f"An error occurred: {e}")
#     finally:
#         browser.quit()


# if __name__ == "__main__":
#     main()



# import os
# import csv
# import re
# import requests
# from bs4 import BeautifulSoup

# BASE_URL = 'https://www.tandem-associates.com/lionel/lionel_trains_master_index.htm#TOPLIST'  # Replace with the actual URL

# response = requests.get(BASE_URL)
# soup = BeautifulSoup(response.content, 'html.parser')

# # Extract all links in the main table
# table_links = soup.select('table a[href]')
# print(f"Found {len(table_links)} links to process.")


# def process_link(tag):
#     link = tag['href']
#     full_url = os.path.join(BASE_URL, link)
#     response = requests.get(full_url)
#     page_soup = BeautifulSoup(response.content, 'html.parser')

#     # Extracting image details
#     img = page_soup.select_one('img[src*="pictures/"]')
#     if img:
#         img_src = os.path.join(BASE_URL, img['src'])
#         img_alt = img.get('alt', '')
#     else:
#         img_src = ''
#         img_alt = ''

#     # Extracting name and model number
#     text_content = page_soup.text
#     name_match = re.search(r'Lionel Trains(.*?)No\. (\d+)', text_content)
#     if name_match:
#         name = name_match.group(1).strip()
#         model_no = name_match.group(2).strip()
#     else:
#         name = ''
#         model_no = ''

#     return name, model_no, img_src, img_alt, link


# data = [process_link(tag) for tag in table_links]

# with open('output.csv', 'w', newline='') as csvfile:
#     writer = csv.writer(csvfile)
#     writer.writerow(["Name", "Model Number", "Image Link", "Alt Text", "Href Name"])
#     writer.writerows(data)

# print("Data saved to output.csv")
