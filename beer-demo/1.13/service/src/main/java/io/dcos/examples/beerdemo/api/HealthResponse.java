package io.dcos.examples.beerdemo.api;

public class HealthResponse {
  public boolean healthy;

  public HealthResponse(boolean healthy) {
    this.healthy = healthy;
  }
}
