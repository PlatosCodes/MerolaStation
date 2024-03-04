import pandas as pd
import tkinter as tk
from tkinter import simpledialog, messagebox

# Load CSV file
csv_path = '../train_data.csv'
df = pd.read_csv(csv_path, header=None, names=['Full', 'Model', 'Name'])

def check_for_unwanted_strings(cell):
    unwanted_strings = ["!", "(", ")", "Mint"]
    for s in unwanted_strings:
        if s in cell:
            return True
    return False

# Check for unwanted strings and ask user for replacements
for idx, row in df.iterrows():
    # Skip 'Full' column, so we start checking from 'Model' and 'Name'
    for col in ['Model', 'Name']:
        if check_for_unwanted_strings(row[col]):
            root = tk.Tk()
            root.withdraw()  # Hide the root window
            new_name = simpledialog.askstring("Input", f"Found unwanted string in {col}: {row[col]}. Enter new name for {col}:")
            df.at[idx, col] = new_name
            root.destroy()

# Check for duplicates
checked_models = set()
for idx, row in df.iterrows():
    model = row['Model']
    if model in checked_models:
        # Prompt user for action
        root = tk.Tk()
        root.withdraw()  # Hide the root window
        choice = simpledialog.askinteger("Choose Entry", f"Duplicate model {model} found.\n1. Keep first\n2. Keep second", minvalue=1, maxvalue=2)
        if choice == 1:
            df = df.drop(idx)
        else:
            # Drop all rows with this model and then add the current row back
            df = df[df['Model'] != model]
            df = df.append(row)
        root.destroy()
    checked_models.add(model)
# Write updated data to CSV
df.to_csv(csv_path+"new", index=False, header=False)
messagebox.showinfo("Info", "CSV file updated successfully!")

