import { Request, Response, NextFunction } from "express";

export interface LoggingOptions {
  logRequests?: boolean;
  logResponses?: boolean;
  logHeaders?: boolean;
}

export function createLoggingMiddleware(options: LoggingOptions = {}) {
  const {
    logRequests = true,
    logResponses = true,
    logHeaders = false
  } = options;

  return (req: Request, res: Response, next: NextFunction) => {
    const start = Date.now();

    if (logRequests) {
      const requestInfo = [
        `[REQUEST]`,
        new Date().toISOString(),
        req.ip || req.connection.remoteAddress,
        req.method,
        req.path
      ];

      if (logHeaders && Object.keys(req.headers).length > 0) {
        requestInfo.push(`Headers: ${JSON.stringify(req.headers)}`);
      }

      console.log(requestInfo.join(" | "));
    }

    // Capture original res.end to log response
    const originalEnd = res.end.bind(res);
    res.end = function(chunk?: any, encoding?: any, callback?: () => void) {
      if (logResponses) {
        const duration = Date.now() - start;
        const responseInfo = [
          `[RESPONSE]`,
          new Date().toISOString(),
          req.ip || req.connection.remoteAddress,
          req.method,
          req.path,
          `Status: ${res.statusCode}`,
          `Duration: ${duration}ms`
        ];

        console.log(responseInfo.join(" | "));
      }

      // Call original end method
      return originalEnd(chunk, encoding, callback);
    };

    next();
  };
}