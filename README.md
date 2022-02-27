# memento

Memento is a command line tool written in go for sorting and categorizing personal files like screenshots, recordings, logs and more. It can generate a report for a specific day, containing all the datapoints from that day, giving the user a timeline-like overview  of what happened that day. 

## Supported platforms

- Windows
  
  Windows is fully supported, including permanent data sources

- Linux and Mac
  
  Linux and Mac is supported, excluding permanent data sources (support for permanent data sources on those platforms is coming soon!)

## Installation

Clone this repo with

```
git clone https://github.com/MergenStudios/memento
```

### Windows

Open a new terminal as administrator (this is required to add the service that manages permanent data sources). Navigate to the folder and run

```
make windows
```

### Mac/Linux

Navigate to the folder and run

```
go install memento.go
```



To use, type `memento` anywhere in your terminal

## Usage

To get started, use `memento setup` to generate the required directory structure. A `typesEnums.json` can be found in the `config` folder. The datatypes in your memento project are stored here. You can use `memento types` to manipulate them. 

To actually import data into your memento project, you can use `memento import`. Lets say you have a folder of old family photos you want to import to memento. For that you would use `memento import PHOTOS Path/To/Photos`.

To generate a report for a specific day, you can use `memento report`. Going back to the previous example, lets say you want to get an overview of the first day of your trip to Berlin. Knowing the date and timezone, you can run `memento report 2018-05-28 Europe/Berlin`, creating `Report-2018-05-28.txt` in the reports folder, which might look something like this:

```
Report 28-05-2018

10:43:35            | Photo        Path/To/Photos/photo231.png
10:54:39            | Photo        Path/To/Photos/photo232.png
15:38:50            | Photo        Path/To/Photos/photo233.png
```

For the full documentation you can always use `memento [command] --help`

**Contributions are welcome, feel free to make a pull request or open an issue!**
