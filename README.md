# Hygieia

Hygieia is the Ancient Egyptian god of cleanliness.  This project helps enable VATSIM Facility Engineers to easily filter out unwanted lines within their diagrams.

## Configuration

The configuration is a YAML file.

```yaml
filter: inside
points:
- lat: 43.434109
  lon: -88.890204
- lat: 44.375058
  lon: -88.628186
- lat: 44.411986
  lon: -89.561127
- lat: 43.425246
  lon: -89.966063
- lat: 43.434109
  lon: -88.890204
```

### Filter
* inside: Any line that has a point within the polygon will be filtered
* outside: Any line that has a point outside of the polygon will be filtered

### Points
* An object of lat and lons that create the polygon to use for filtering

## Usage

To run (using Linux as an example), run the following. In this example I am using Linux and will be cleaning "Chicago ARTCC Combined_1617077493.sct2" and want to output to cleaned.sct2.

```
./hygieia_linux_amd64 clean Chicago\ ARTCC\ Combined_1617077493.sct2 output.sct2
```

The command syntax is the same, replace the executable as necessary for your OS.

Run the command help for more information.
```
./hygieia_linux_amd64 help
```

## License

This project is licensed under GPL 3. Please see LICENSE.md for more information.