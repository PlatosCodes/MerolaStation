from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service as ChromeService
from fake_useragent import UserAgent
import time
import random
import csv, string
import tkinter as tk
from tkinter import simpledialog

# This is Google's base URL for image searches.
GOOGLE_IMG_SEARCH_URL = 'https://www.google.com/search?tbm=isch&q='

def prompt_for_new_model(original_model):
    """Prompt the user for a new model number or skip action."""
    root = tk.Tk()
    root.withdraw()  # hide the main window
    user_input = simpledialog.askstring("Duplicate Entry Detected", f"Duplicate model detected for: {original_model}. \nProvide a new model number or type 'skip'.")
    root.destroy()
    return user_input

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


# class GUIHandler:
#     def __init__(self):
#         self.root = tk.Tk()
#         self.root.withdraw()

#     def prompt_for_new_model(self, original_model):
#         user_input = simpledialog.askstring("Duplicate Entry Detected", f"Duplicate model detected for: {original_model}. \nProvide a new model number or type 'skip'.", parent=self.root)
#         return user_input

#     def destroy(self):
#         self.root.destroy()



def get_train_image_url(browser, train_name):
    """Extract the first image URL from Google Image Search for a given train name."""
    search_url = GOOGLE_IMG_SEARCH_URL + train_name
    browser.get(search_url)
    time.sleep(random.uniform(1, 2))  # Random sleep time between 1 to 3 seconds to mimic human behavior

    # Get the source URL of the first image result
    try:
        img_element = browser.find_element(By.CSS_SELECTOR, 'img.rg_i')
        img_url = img_element.get_attribute('src')
        return img_url
    except Exception as e:
        print(f"Error extracting image for {train_name}. Error: {e}")
        return None

def main():
    browser = setup_browser()
    # gui = GUIHandler()

    merged_data = []
    model_numbers = set()

    # Open the first CSV, read the data, and add to merged_data
    with open('../train_data.csv', 'r') as csvfile:
        reader = csv.reader(csvfile)
        for row in reader:
            model_numbers.add(row[1])  # Assuming the model number is in the second column (index 1)
            merged_data.append(row)

    with open('./merged1.csv', 'r') as csvfile:
        reader = csv.reader(csvfile)
        for row in reader:
            model_numbers.add(row[1])  # Assuming the model number is in the second column (index 1)
            merged_data.append(row)

    i = 0
    # Open the second CSV, check for duplicates, fetch image URLs, and merge data
    with open('./train_data.csv', 'r') as csvfile:
        reader = csv.reader(csvfile)
        for row in reader:
            if row[1] in model_numbers:
                row[1] = row[1]+"second"
            train_name = row[0]
            img_url = get_train_image_url(browser, train_name)
            i += 1
            if img_url:
                merged_data.append(row + [img_url])
                print("got img for:", train_name, img_url, "image #" + str(i))
            else:
                merged_data.append(row + [""])

    # Save the merged data to a new CSV
    with open('final_merge.csv', 'w', newline='') as csvfile:
        writer = csv.writer(csvfile)
        for row in merged_data:
            writer.writerow(row)

    browser.quit()
    # gui.destroy()

if __name__ == "__main__":
    main()