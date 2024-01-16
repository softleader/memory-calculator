package org.example.jibdemo;

import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.ApplicationRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;

@Slf4j
@SpringBootApplication
public class JibDemoApplication {

  public static void main(String[] args) {
    SpringApplication.run(JibDemoApplication.class, args);
  }

  @Bean
  ApplicationRunner printEnv() {
    return args -> {
      log.info("=== System Env ===");
      System.getenv().forEach((k, v) -> log.info("{}={}", k, v));

      log.info("=== System Properties ===");
      System.getProperties().forEach((k, v) -> log.info("{}={}", k, v));

      log.info("=== Received {} Args ===", args.getSourceArgs().length);
      for (String arg : args.getSourceArgs()) {
        log.info("{}", arg);
      }
    };
  }
}
