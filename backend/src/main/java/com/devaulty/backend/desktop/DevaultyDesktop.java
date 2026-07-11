package com.devaulty.backend.desktop;

import javafx.application.Application;
import javafx.scene.Scene;
import javafx.scene.layout.StackPane;
import javafx.scene.web.WebView;
import javafx.stage.Stage;

/**
 * Entry point for Devaulty in desktop mode (JavaFX + WebView).
 *
 * WARNING: This class is still under development / proof of concept.
 * Currently, the WebView points to a frontend server running separately
 * (Vite dev/preview on localhost), used only to validate the JavaFX +
 * React integration during development.
 *
 * This does NOT represent the final packaged version of the app:
 *   - The frontend is not yet embedded/served alongside the backend.
 *   - The Spring Boot backend is not automatically started from here.
 *   - There is no error handling for connection failures, splash
 *     screen, or any production-level UX.
 *
 */

public class DevaultyDesktop extends Application {

    private static final String FRONTEND_DEV_URL = "http://localhost:4173";

    @Override
    public void start(Stage primaryStage) {

        WebView webView = new WebView();
        webView.getEngine().load(FRONTEND_DEV_URL);

        StackPane root = new StackPane(webView);
        Scene scene = new Scene(root, 1200, 800);

        primaryStage.setTitle("Devaulty Desktop");
        primaryStage.setScene(scene);
        primaryStage.show();
    }

    public static void main(String[] args) {
        launch(args);
    }
}