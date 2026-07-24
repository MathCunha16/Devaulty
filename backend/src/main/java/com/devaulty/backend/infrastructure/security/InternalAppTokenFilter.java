package com.devaulty.backend.infrastructure.security;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;

@Component
@Order(1)
public class InternalAppTokenFilter extends OncePerRequestFilter {

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

        if (requestToken == null || !requestToken.equals(AppTokenContext.PROCESS_TOKEN)) {
            response.setStatus(HttpServletResponse.SC_FORBIDDEN);
            response.setContentType("application/json");
            response.getWriter().write("{\"error\": \"Forbidden: Unauthorized local process request.\"}");
            return;
        }

        filterChain.doFilter(request, response);
    }
}
