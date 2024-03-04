import requests
import csv

EBAY_ACCESS_TOKEN = "YOUR_EBAY_ACCESS_TOKEN"
CSV_FILE_PATH = "path_to_your_csv_file.csv"

def fetch_average_price_from_ebay(query):
    endpoint = "https://api.ebay.com/buy/browse/v1/item_summary/search_by_keyword"
    headers = {
        "Authorization": f"Bearer {EBAY_ACCESS_TOKEN}",
        "X-EBAY-C-ENDUSERCTX": "contextualLocation=country=US,zip=ZIP_CODE"
    }
    params = {
        "q": query,
        "sort": "price(asc)",
        "filter": "itemLocationCountry:US"
    }
    
    response = requests.get(endpoint, headers=headers, params=params)
    data = response.json()
    
    if 'itemSummaries' not in data:
        return 0

    total_value = sum([item['price']['value'] for item in data['itemSummaries']])
    average_value = total_value / len(data['itemSummaries'])
    return average_value

def process_csv_and_fetch_prices():
    with open(CSV_FILE_PATH, 'r') as file:
        csv_reader = csv.reader(file)
        for row in csv_reader:
            query = row[0]
            average_price = fetch_average_price_from_ebay(query)
            print(f"Average price for {query}: ${average_price:.2f}")

if __name__ == "__main__":
    process_csv_and_fetch_prices()
