package com.devaulty.backend.desktop;

import com.devaulty.backend.BackendApplication;
import com.devaulty.backend.desktop.listener.ServerPortListener;
import com.devaulty.backend.infrastructure.security.AppTokenContext;
import javafx.application.Application;
import javafx.application.Platform;
import javafx.geometry.Rectangle2D;
import javafx.scene.Scene;
import javafx.scene.control.Alert;
import javafx.scene.image.Image;
import javafx.scene.image.ImageView;
import javafx.scene.paint.Color;
import javafx.scene.web.WebView;
import javafx.stage.Screen;
import javafx.stage.Stage;
import javafx.stage.StageStyle;
import org.springframework.boot.builder.SpringApplicationBuilder;
import org.springframework.context.ConfigurableApplicationContext;

import java.io.InputStream;

public class DevaultyDesktop extends Application {

    private static final double SPLASH_WIDTH = 550;
    private static final double SPLASH_HEIGHT = 400;

    private ConfigurableApplicationContext springContext;
    private Stage splashStage;

    @Override
    public void init() {
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

        // 1. Create and show the splash screen immediately
        createAndShowSplashScreen();

        // 2. Prepares the main WebView window but does not display yet
        WebView webView = new WebView();
        webView.setContextMenuEnabled(false);

        Screen screen = Screen.getPrimary();
        Rectangle2D bounds = screen.getVisualBounds();

        // Window is sized/positioned to cover the primary screen entirely,
        // simulating a maximized state manually instead of calling
        // setMaximized(true). On several Linux/GTK window managers,
        // setMaximized(true) decides which monitor to maximize onto based on
        // cursor position / "active" monitor, ignoring the Stage's X/Y — which
        // is what was causing the window to jump to the wrong screen.
        Scene mainScene = new Scene(webView, bounds.getWidth(), bounds.getHeight());
        primaryStage.setTitle("Devaulty");
        primaryStage.setScene(mainScene);

        InputStream iconStream = getClass().getResourceAsStream("/static/devaulty-icon.png");
        if (iconStream != null) {
            primaryStage.getIcons().add(new Image(iconStream));
        }

        primaryStage.setMinWidth(1000);
        primaryStage.setMinHeight(650);

        primaryStage.setX(bounds.getMinX());
        primaryStage.setY(bounds.getMinY());
        primaryStage.setWidth(bounds.getWidth());
        primaryStage.setHeight(bounds.getHeight());

        // 3. Waits for the server port to be available
        ServerPortListener.serverPortFuture.thenAccept(port -> {
            String appUrl = "http://localhost:" + port;

            Platform.runLater(() -> {
                webView.getEngine().load(appUrl);

                // When HTML is loaded, injects the internal token
                webView.getEngine().getLoadWorker().stateProperty().addListener((obs, oldState, newState) -> {
                    if (newState == javafx.concurrent.Worker.State.SUCCEEDED) {

                        // Injects the internal token into the WebView
                        String script = String.format("window.DEVAULTY_INTERNAL_TOKEN = '%s';", AppTokenContext.PROCESS_TOKEN);
                        webView.getEngine().executeScript(script);

                        primaryStage.show();
                        // No setMaximized(true) here — size/position were
                        // already set above to cover the full primary screen.

                        if (splashStage != null) {
                            splashStage.close();
                        }
                    }
                });
            });
        }).exceptionally(ex -> {
            Platform.runLater(() -> showErrorAndExit("Failed to obtain server port", ex));
            return null;
        });
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

        javafx.scene.layout.StackPane root = new javafx.scene.layout.StackPane();
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