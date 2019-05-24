package io.dcos.examples.beerdemo.api;

public class BeerResponse {
  public String hostAddress;
  public String uuid;
  public String version;
  public String beerName;
  public String beerStyle;
  public String beerDescription;

  public BeerResponse(String hostAddress, String uuid, String version, String beerName, String beerStyle, String beerDescription) {
    this.hostAddress = hostAddress;
    this.uuid = uuid;
    this.version = version;
    this.beerName = beerName;
    this.beerStyle = beerStyle;
    this.beerDescription = beerDescription;
  }
}
