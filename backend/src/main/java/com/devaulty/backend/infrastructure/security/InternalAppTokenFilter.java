package com.devaulty.backend.infrastructure.security;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;

@Component
@Order(1)
public class InternalAppTokenFilter extends OncePerRequestFilter {

    @Value("${devaulty.dev.token:#{null}}")
    private String devToken;

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain) throws ServletException, IOException {

        String path = request.getRequestURI();

        // 1. Allows static resources from React, Swagger and OPTIONS (Pre-flight CORS on dev)
        if (!path.startsWith("/api/") || "OPTIONS".equalsIgnoreCase(request.getMethod())) {
            filterChain.doFilter(request, response);
            return;
        }

        // 2. For API requests, requires the request token
        String requestToken = request.getHeader(AppTokenContext.HEADER_NAME);

        // 3. Safety check
        // In Dev: accepts the devToken
        // In Prod: accepts the processToken
        boolean isDevTokenValid = devToken != null && !devToken.isBlank() && devToken.equals(requestToken);
        boolean isProcessTokenValid = AppTokenContext.PROCESS_TOKEN.equals(requestToken);

        if (requestToken == null || (!isDevTokenValid && !isProcessTokenValid)) {
            response.setStatus(HttpServletResponse.SC_FORBIDDEN);
            response.setContentType("application/json");
            response.getWriter().write("{\"error\": \"Forbidden: Unauthorized local process request.\"}");
            return;
        }

        filterChain.doFilter(request, response);
    }
}
