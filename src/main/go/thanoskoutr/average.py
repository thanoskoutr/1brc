#!/bin/python3

def convert_to_seconds(time_str):
    if ':' in time_str:
        minutes, seconds = time_str.split(':')
        seconds = seconds.replace(',', '.')
        return int(minutes) * 60 + float(seconds)
    else:
        return float(time_str.replace(',', '.'))

def convert_to_time_format(total_seconds):
    minutes = int(total_seconds // 60)
    seconds = total_seconds % 60
    return f"{minutes}:{seconds:05.2f}"

# Execution times (Fill the execution times here)
time1 = "4,701"
time2 = "4,644"
time3 = "4,630"

# Convert times to seconds
seconds1 = convert_to_seconds(time1)
seconds2 = convert_to_seconds(time2)
seconds3 = convert_to_seconds(time3)

# Calculate average
average_seconds = (seconds1 + seconds2 + seconds3) / 3
print("Average execution time (seconds):", average_seconds)

# Convert average back to m:s.ms format
average_time = convert_to_time_format(average_seconds)

print(f"Average execution time: {average_time}")