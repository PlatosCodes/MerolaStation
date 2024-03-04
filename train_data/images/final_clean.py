import csv

def process_csv(input_file, output_file):
    with open(input_file, 'r') as infile, open(output_file, 'w', newline='') as outfile:
        reader = csv.reader(infile)
        writer = csv.writer(outfile)

        for row in reader:
            # Check if the model number is "0"
            if row[1] == "0":
                continue

            # Ensure the first column is the concatenation of the second and third columns
            row[0] = row[1] + " " + row[2]

            # Write the corrected/verified row to the new CSV
            writer.writerow(row)

input_file = "train_data_cleaned.csv"
output_file = "train_data.csv"
process_csv(input_file, output_file)
