import csv
import re

# Define a function to extract the src value from the img tag
def extract_src(img_tag):
    match = re.search(r'src="(.*?)"', img_tag)
    return match.group(1) if match else None

# Define a function to check if the row matches the criteria
def is_desired_img_tag(row):
    # Regular expression to match the criteria
    pattern = r'src=".*_ident\.(gif|jpg)"'
    return re.search(pattern, row[0])

# Set to hold the unique src values
seen_srcs = set()

# List to hold the filtered rows
filtered_rows = []

# Open the CSV and read its contents
with open('train_images.csv', 'r') as csvfile:
    reader = csv.reader(csvfile)
    for row in reader:
        src_value = extract_src(row[0])
        if is_desired_img_tag(row) and src_value not in seen_srcs:
            seen_srcs.add(src_value)
            filtered_rows.append(row)

# Write the filtered rows to a new CSV
with open('cleaned_train_images.csv', 'w', newline='') as csvfile:
    writer = csv.writer(csvfile)
    writer.writerows(filtered_rows)
