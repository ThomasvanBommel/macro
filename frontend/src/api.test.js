import { test, expect } from "vitest";
import { http, HttpResponse } from "msw";
import { setupServer } from "msw/node";

import { postJSON } from "./api";

const api = path => new URL(path, "http://localhost").href;

export const handlers = [
    http.post(api("/test"), (req) => {
        const { name } = req.json();
        if (name === "error") {
            return HttpResponse.json({ message: "Test error" }, { status: 400 });
        }
        return HttpResponse.json({ message: `Hello, ${name}!` }, { status: 200 });
    }),
];

export const server = setupServer(...handlers);

test("postJSON returns data on success", async () => {
    const data = await postJSON(api("/test"), { name: "Alice" });
    expect(data).resolves.toEqual({ message: "Hello, Alice!" });
});

test("postJSON throws an error on failure", async () => {
    try {
        await postJSON(api("/test"), { name: "error" });
        throw new Error("Expected postJSON to throw an error");
    } catch (err) {
        expect(err).toBeInstanceOf(Error);
        expect(err.message).resolves.toEqual({ message: "Test error" });
        expect(err.status).toEqual(400);
    }
});