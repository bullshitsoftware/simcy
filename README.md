# Simcy

An easy way to keep your [SimulationCraft.exe](https://www.simulationcraft.org/) updated.

## Usage 

```
> .\simcy.exe --help
Usage of D:\simcy.exe:
  -downloads-url string
        simcraft downloads page (default "http://downloads.simulationcraft.org/nightly/")
  -storage string
        where to store simcraft (default "C:\\Users\\username\\AppData\\Local/Simcy")
```

The tool works as follows:

1. Checks the last version at downloads page
2. If you don't have last version, tool will download and extract it to the `-storage` location. 
After this tool will purge all other directories at the storage (be careful with custom path)
3. Then the latest SimulationCraft will be run
