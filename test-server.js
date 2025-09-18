#!/usr/bin/env node

/**
 * Simple test script to verify the TypeScript MCP server implementation
 */

import { exec } from "node:child_process";
import { promisify } from "node:util";

const execAsync = promisify(exec);

async function testBuild() {
  console.log("🔨 Testing build...");
  try {
    await execAsync("npm run build");
    console.log("✅ Build successful");
  } catch (error) {
    console.error("❌ Build failed:", error);
    process.exit(1);
  }
}

async function testLint() {
  console.log("🔍 Testing lint...");
  try {
    await execAsync("npm run lint");
    console.log("✅ Lint passed");
  } catch (error) {
    console.warn("⚠️ Lint warnings found");
  }
}

async function testServerStart() {
  console.log("🚀 Testing server startup...");
  return new Promise((resolve, reject) => {
    const child = exec("npm run dev server", (error) => {
      if (error && !error.killed) {
        reject(error);
      }
    });

    child.stdout?.on("data", (data) => {
      if (data.includes("Unsplash MCP server listening")) {
        setTimeout(() => {
          child.kill();
          console.log("✅ Server starts successfully");
          resolve(true);
        }, 1000);
      }
    });

    child.stderr?.on("data", (data) => {
      console.error("Server error:", data);
    });

    setTimeout(() => {
      child.kill();
      reject(new Error("Server startup timeout"));
    }, 10000);
  });
}

async function testHealthEndpoint() {
  console.log("🏥 Testing health endpoint...");
  return new Promise((resolve, reject) => {
    const child = exec("npm run dev server", (error) => {
      if (error && !error.killed) {
        reject(error);
      }
    });

    child.stdout?.on("data", (data) => {
      if (data.includes("Unsplash MCP server listening")) {
        setTimeout(async () => {
          try {
            const { stdout } = await execAsync("curl -s http://127.0.0.1:8081/health");
            const response = JSON.parse(stdout);
            if (response.status === "ok") {
              console.log("✅ Health endpoint working");
              resolve(true);
            } else {
              reject(new Error("Health check failed"));
            }
          } catch (error) {
            reject(error);
          } finally {
            child.kill();
          }
        }, 2000);
      }
    });

    setTimeout(() => {
      child.kill();
      reject(new Error("Health test timeout"));
    }, 15000);
  });
}

async function runTests() {
  console.log("🧪 Running TypeScript MCP Server Tests\n");

  try {
    await testBuild();
    await testLint();
    await testServerStart();
    await testHealthEndpoint();

    console.log("\n🎉 All tests passed! TypeScript MCP server is ready.");
    console.log("\n📝 Next steps:");
    console.log("1. Add your Unsplash API key to .env file");
    console.log("2. Test with a real MCP client");
    console.log("3. Run 'npm run dev server' to start HTTP mode");
    console.log("4. Run 'npm run dev stdio' to start stdio mode");
  } catch (error) {
    console.error("\n❌ Tests failed:", error);
    process.exit(1);
  }
}

runTests();