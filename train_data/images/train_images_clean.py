import csv
import re
import tkinter as tk
from tkinter import simpledialog

def extract_data(row):
    # Extracting the title and image src using regex
    title_match = re.search(r'title="([^"]+)"', row)
    src_match = re.search(r'src="([^"]+)"', row)

    if title_match and src_match:
        title = title_match.group(1)
        
        # Attempt to extract model number from the two possible structures
        model_number_match = re.search(r'No\. ([^\s]+)', title)
        model_number = model_number_match.group(1) if model_number_match else ""
        
        # Depending on the location of the model number, determine the name
        if title.startswith("No."):
            name = title.split(model_number, 1)[-1].strip()
        else:
            name = title.split("No.")[0].strip()

        # Remove undesired words from name
        name = name.replace("Lionel Trains", "").replace("The Lionel", "").strip()

        src = src_match.group(1)
        
        return [model_number + " " + name, model_number, name, src]
    return None


def get_input(prompt):
    root = tk.Tk()
    root.withdraw()  # Hide the main window
    user_input = simpledialog.askstring("Input", prompt)
    root.destroy()  # Close the window
    return user_input

def process_csv(input_file, output_file):
    with open(input_file, 'r') as infile, open(output_file, 'w', newline='') as outfile:
        reader = csv.reader(infile)
        writer = csv.writer(outfile)

        for row in reader:
            data = extract_data(row[0])
            if data:
                # Check for empty model number or title
                if not data[1]:
                    data[1] = get_input("Enter model number for title: " + data[2])
                    data[0] = data[1] + " " + data[2]  # Update the concatenated column as well

                if not data[2]:
                    data[2] = get_input("Enter title for model number: " + data[1])
                    data[0] = data[1] + " " + data[2]  # Update the concatenated column as well

                writer.writerow(data)

input_file = "train_images.csv"
output_file = "train_data.csv"
process_csv(input_file, output_file)

                # import csv
# import re

# def extract_data(row):
#     # Extracting the title and image src using regex
#     title_match = re.search(r'title="([^"]+)"', row)
#     src_match = re.search(r'src="([^"]+)"', row)

#     if title_match and src_match:
#         title = title_match.group(1)
        
#         # Attempt to extract model number from the two possible structures
#         model_number_match = re.search(r'No\. ([^\s]+)', title)
#         model_number = model_number_match.group(1) if model_number_match else ""
        
#         # Depending on the location of the model number, determine the name
#         if title.startswith("No."):
#             name = title.split(model_number, 1)[-1].strip()
#         else:
#             name = title.split("No.")[0].strip()

#         # Remove undesired words from name
#         name = name.replace("Lionel Trains", "").replace("The Lionel", "").strip()

#         src = src_match.group(1)
        
#         return [model_number + " " + name, model_number, name, src]
#     return None


# def process_csv(input_file, output_file):
#     with open(input_file, 'r') as infile, open(output_file, 'w', newline='') as outfile:
#         reader = csv.reader(infile)
#         writer = csv.writer(outfile)

#         for row in reader:
#             data = extract_data(row[0])
#             if data:
#                 writer.writerow(data)

# input_file = "train_images.csv"
# output_file = "./train_data.csv"
# process_csv(input_file, output_file)
