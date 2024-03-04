import csv

def process_csv(input_file, output_file):
    base_url = "https://www.tandem-associates.com/lionel/"

    with open(input_file, 'r') as infile, open(output_file, 'w', newline='') as outfile:
        reader = csv.reader(infile)
        writer = csv.writer(outfile)

        for row in reader:
            # Update the fourth column
            row[3] = base_url + row[3]

            # Write the modified row to the new CSV
            writer.writerow(row)

input_file = "train_data_without_main_url.csv"
output_file = "train_data.csv"
process_csv(input_file, output_file)
