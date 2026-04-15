import { describe, test, expect } from "vitest";
import { http, HttpResponse } from "msw";
import { setupServer } from "msw/node";

import { postJSON } from "./api";

const api = path => new URL(path, "http://localhost").href;

export const server = setupServer(
    http.post(api("/test"), async r => {
        const { name } = await r.request.json();

        switch (name) {
            case "error":
                return HttpResponse.json({ error: "Test error" }, { status: 400 });
            case "unknown":
                return HttpResponse.json({}, { status: 400 });
            default:
                return HttpResponse.json({ message: `Hello, ${name}!` }, { status: 200 });
        }
    })
);

describe("postJSON", () => {
    test("returns data on success", async () => {
        const data = await postJSON(api("/test"), { name: "Alice" });
        expect(data).toEqual({ message: "Hello, Alice!" });
    });

    test("throws an error on failure", async () => {
        try {
            await postJSON(api("/test"), { name: "error" });
            throw new Error("Expected postJSON to throw an error");
        } catch (err) {
            expect(err).toBeInstanceOf(Error);
            expect(err.message).toEqual("Test error");
            expect(err.status).toEqual(400);
        }
    });

    test("Unknown bad response", async () => {
        try {
            await postJSON(api("/test"), { name: "unknown" });
            throw new Error("Expected postJSON to throw an error");
        } catch (err) {
            expect(err).toBeInstanceOf(Error);
            expect(err.message).toEqual("Unknown error");
            expect(err.status).toEqual(400);
        }
    });
});