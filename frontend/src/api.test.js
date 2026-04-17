import { beforeAll, afterEach, afterAll, describe, test, expect } from 'vitest'
import { http, HttpResponse } from "msw";
import { setupServer } from "msw/node";

import { postJSON } from "./api";

const server = setupServer(
    http.post("/test", async r => {
        const { name } = await r.request.json();

        switch (name) {
            case "error":
                return HttpResponse.json({ error: "Test error" }, { status: 400 });
            case "unknown":
                return HttpResponse.json({}, { status: 400 });
            default:
                return HttpResponse.json({ message: `Hello, ${name}!` }, { status: 200 });
        }
    }),
);

beforeAll(() => server.listen({ onUnhandledRequest: "error" }));
afterEach(() => server.resetHandlers())
afterAll(()  => server.close())

describe("postJSON", () => {
    test.concurrent("returns data on success", async () => {
        const data = await postJSON("/test", { name: "Alice" });
        expect(data).toEqual({ message: "Hello, Alice!" });
    });

    test.concurrent("throws an error on failure", async () => {
        try {
            await postJSON("/test", { name: "error" });
            throw new Error("Expected postJSON to throw an error");
        } catch (err) {
            expect(err).toBeInstanceOf(Error);
            expect(err.message).toEqual("Test error");
            expect(err.status).toEqual(400);
        }
    });

    test.concurrent("Unknown bad response", async () => {
        try {
            await postJSON("/test", { name: "unknown" });
            throw new Error("Expected postJSON to throw an error");
        } catch (err) {
            expect(err).toBeInstanceOf(Error);
            expect(err.message).toEqual("Unknown error");
            expect(err.status).toEqual(400);
        }
    });
});