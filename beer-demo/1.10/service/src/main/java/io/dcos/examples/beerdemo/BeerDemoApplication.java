package io.dcos.examples.beerdemo;

import io.dcos.examples.beerdemo.api.BeerResponse;
import io.dcos.examples.beerdemo.api.HealthResponse;
import org.apache.commons.io.IOUtils;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.log4j.MDC;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.actuate.health.Health;
import org.springframework.boot.actuate.health.HealthIndicator;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurerAdapter;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Locale;
import java.util.Map;
import java.util.UUID;

@SpringBootApplication
@RestController("/")
public class BeerDemoApplication extends WebMvcConfigurerAdapter {

  // stores the host address on which this service runs
  private final String hostAddress;

  // stores the application version of this service
  @Value("${SERVICE_VERSION:1}")
  private String version;

  @Value("${ELASTICSEARCH_URL:http://localhost:9200}")
  private String elasticsearchUrl;

  private final CloseableHttpClient client = HttpClients.createDefault();

  private String uuid = UUID.randomUUID().toString().replace("-", "");

  // stores the information if this service should be return healthy or unhealthy
  private boolean sober = true;

  private final JdbcTemplate template;

  public BeerDemoApplication(final JdbcTemplate jdbcTemplate) {
    this.template = jdbcTemplate;
    hostAddress = getHostAddress();
    initMDC();
  }

  @Bean
  public HealthIndicator chucksHealthIndicator() {
    return () -> (sober ? Health.up() : Health.down()).build();
  }

  @RequestMapping("/")
  public BeerResponse randomBeer(Locale locale) {
    Map<String, Object> query = template.queryForMap(
        "SELECT b.id, b.name, b.descript, s.style_name " +
            "FROM `beers` b LEFT JOIN `styles` s " +
            "ON b.style_id = s.id " +
            "WHERE descript != '' " +
            "ORDER BY RAND() LIMIT 0,1;");

    return new BeerResponse(
        hostAddress,
        uuid,
        version,
        "" + query.get("name"),
        "" + query.get("style_name"),
        query.get("descript").toString()
    );
  }

  @RequestMapping(value = "/search", produces = "application/json")
  public String searchProxy(@RequestParam("q") String query) {
    try (CloseableHttpResponse response = client.execute(new HttpGet(elasticsearchUrl + "/beer/_search?q=" + query))) {
      return IOUtils.toString(response.getEntity().getContent(), "UTF-8");
    } catch (Exception e) {
      throw new RuntimeException("o_O - Error while proxying elasticsearch", e);
    }
  }

  @RequestMapping(value = "/health", method = RequestMethod.DELETE)
  public HealthResponse toggleHealth() {
    sober = false;
    return new HealthResponse(false);
  }


  public static void main(String[] args) throws InterruptedException {
    if ("true".equals(System.getenv("WAIT_AT_STARTUP"))) {
      Thread.sleep(15000);
    }
    SpringApplication.run(BeerDemoApplication.class, args);
  }

  private String getHostAddress() {
    try {
      return InetAddress.getLocalHost().getHostAddress();
    } catch (UnknownHostException e) {
      return "unknown";
    }
  }

  private void initMDC() {
    MDC.put("uuid", uuid);
    MDC.put("version", version);
    MDC.put("host", hostAddress);
  }
}
