package com.devaulty.backend.desktop;

import com.devaulty.backend.BackendApplication;
import com.devaulty.backend.desktop.listener.ServerPortListener;
import com.devaulty.backend.infrastructure.security.AppTokenContext;
import javafx.application.Application;
import javafx.application.Platform;
import javafx.concurrent.Worker;
import javafx.geometry.Rectangle2D;
import javafx.scene.Scene;
import javafx.scene.control.Alert;
import javafx.scene.image.Image;
import javafx.scene.image.ImageView;
import javafx.scene.layout.StackPane;
import javafx.scene.paint.Color;
import javafx.scene.web.WebView;
import javafx.stage.Screen;
import javafx.stage.Stage;
import javafx.stage.StageStyle;
import org.springframework.boot.builder.SpringApplicationBuilder;
import org.springframework.context.ConfigurableApplicationContext;

import java.io.File;
import java.io.InputStream;
import java.net.URL;
import java.util.concurrent.TimeUnit;

public class DevaultyDesktop extends Application {

    private static final double SPLASH_WIDTH = 550;
    private static final double SPLASH_HEIGHT = 400;

    private ConfigurableApplicationContext springContext;
    private Stage splashStage;

    @Override
    public void init() {
        // Ensures user configuration directory exists before Spring Boot initializes SQLite
        File devaultyDir = new File(System.getProperty("user.home"), ".config/devaulty");
        if (!devaultyDir.exists()) {
            devaultyDir.mkdirs();
        }

        // Starts Spring Boot in background without blocking JavaFX UI Thread
        Thread springThread = new Thread(() -> {
            try {
                this.springContext = new SpringApplicationBuilder(BackendApplication.class)
                        .properties("server.port=0")
                        .run();
            } catch (Exception e) {
                Platform.runLater(() -> showErrorAndExit("Failed to start Devaulty Local Server", e));
            }
        }, "devaulty-spring-boot");

        springThread.setDaemon(true); // Closes JVM when Spring Boot is finished
        springThread.start();
    }

    @Override
    public void start(Stage primaryStage) {
        createAndShowSplashScreen();

        WebView webView = new WebView();
        webView.setContextMenuEnabled(false);

        Rectangle2D bounds = Screen.getPrimary().getVisualBounds();
        configurePrimaryStage(primaryStage, webView, bounds);

        ServerPortListener.serverPortFuture
                .orTimeout(60, TimeUnit.SECONDS)
                .thenAccept(port -> Platform.runLater(() -> loadApplication(webView, primaryStage, port)))
                .exceptionally(ex -> {
                    Platform.runLater(() -> showErrorAndExit("Failed to obtain server port", ex));
                    return null;
                });
    }

    private void configurePrimaryStage(Stage primaryStage, WebView webView, Rectangle2D bounds) {
        Scene mainScene = new Scene(webView, bounds.getWidth(), bounds.getHeight());
        primaryStage.setTitle("Devaulty");
        primaryStage.setScene(mainScene);

        URL iconUrl = getClass().getResource("/static/icon/devaulty-icon.png");
        if (iconUrl != null) {
            primaryStage.getIcons().add(new Image(iconUrl.toExternalForm()));
        }

        primaryStage.setMinWidth(1000);
        primaryStage.setMinHeight(650);
        primaryStage.setX(bounds.getMinX());
        primaryStage.setY(bounds.getMinY());
        primaryStage.setWidth(bounds.getWidth());
        primaryStage.setHeight(bounds.getHeight());
    }

    private void loadApplication(WebView webView, Stage primaryStage, Integer port) {
        String appUrl = "http://localhost:" + port;
        webView.getEngine().load(appUrl);

        webView.getEngine().getLoadWorker().stateProperty()
                .addListener((obs, oldState, newState) -> handleWorkerStateChange(webView, primaryStage, newState));
    }

    private void handleWorkerStateChange(WebView webView, Stage primaryStage, Worker.State newState) {
        if (newState == Worker.State.RUNNING) {
            injectInternalToken(webView);
        } else if (newState == Worker.State.SUCCEEDED) {
            injectInternalToken(webView);
            showMainStageAndCloseSplash(primaryStage);
        } else if(newState == Worker.State.FAILED || newState == Worker.State.CANCELLED) {
            showErrorAndExit("Failed to load Devaulty Interface",
                    webView.getEngine().getLoadWorker().getException());
        }
    }

    private void injectInternalToken(WebView webView) {
        try {
            String script = String.format("window.DEVAULTY_INTERNAL_TOKEN = '%s';", AppTokenContext.PROCESS_TOKEN);
            webView.getEngine().executeScript(script);
        } catch (Exception ignored) {
            // Safe to ignore if JavaScript execution context is not yet available or already disposed
        }
    }

    private void showMainStageAndCloseSplash(Stage primaryStage) {
        primaryStage.show();
        if (splashStage != null) {
            splashStage.close();
        }
    }

    @Override
    public void stop() {
        if (this.springContext != null) {
            this.springContext.close(); // Closes SQLite and kills Tomcat gracefully
        }
    }

    private void createAndShowSplashScreen() {
        splashStage = new Stage();
        splashStage.initStyle(StageStyle.TRANSPARENT);

        URL iconUrl = getClass().getResource("/static/icon/devaulty-icon.png");
        if (iconUrl != null) {
            splashStage.getIcons().add(new Image(iconUrl.toExternalForm()));
        }

        StackPane root = new StackPane();
        root.setStyle("-fx-background-color: transparent;");

        InputStream logoStream = getClass().getResourceAsStream("/static/devaulty-splash-screen-logo.png");
        if (logoStream != null) {
            ImageView logoView = new ImageView(new Image(logoStream));
            logoView.setFitWidth(550);
            logoView.setPreserveRatio(true);
            logoView.setSmooth(true);
            root.getChildren().add(logoView);
        }

        Scene splashScene = new Scene(root, SPLASH_WIDTH, SPLASH_HEIGHT);
        splashScene.setFill(Color.TRANSPARENT);

        splashStage.setTitle("Devaulty");
        splashStage.setScene(splashScene);

        Rectangle2D bounds = Screen.getPrimary().getVisualBounds();
        splashStage.setX(bounds.getMinX() + (bounds.getWidth() - SPLASH_WIDTH) / 2);
        splashStage.setY(bounds.getMinY() + (bounds.getHeight() - SPLASH_HEIGHT) / 2);

        splashStage.show();
    }

    private void showErrorAndExit(String message, Throwable cause) {
        Alert alert = new Alert(Alert.AlertType.ERROR);
        alert.setTitle("Initialization Error");
        alert.setHeaderText(message);
        alert.setContentText(cause.getLocalizedMessage());
        alert.showAndWait();
        Platform.exit();
        System.exit(1);
    }
}