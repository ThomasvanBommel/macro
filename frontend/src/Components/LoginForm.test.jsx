import { beforeAll, beforeEach, afterEach, afterAll, describe, test, expect, vi } from 'vitest'
import { render, waitFor, fireEvent, cleanup } from "@testing-library/react";
import { http, HttpResponse } from "msw";
import { setupServer } from "msw/node";

import LoginForm from "./LoginForm";

const server = setupServer(
    http.post("/api/login", async r => {
        const { name } = await r.request.json();

        if (name === "error")
            return HttpResponse.json({ error: "Server error" }, { status: 500 });

        return HttpResponse.json(null, { status: 200 });
    }),
);

function setup() {
    const onError = vi.fn();
    const onSuccess = vi.fn();

    const r = render(<LoginForm onError={onError} onSuccess={onSuccess} />);

    const input = {
        name:     r.getByLabelText(/^username/i),
        password: r.getByLabelText(/^password/i),
        submit:   r.getByRole("button", { type: "submit" })
    };

    return { input, onError, onSuccess };
}

beforeAll(() => server.listen({ onUnhandledRequest: "error" }));
beforeEach(() => server.resetHandlers());
afterEach(() => cleanup());
afterAll(() => server.close());

describe("LoginForm", () => {
    const fast = callback => waitFor(callback, { interval: 2 });

    test("initial state", () => {
        const { input } = setup();

        expect(input.name.value).toBe("");
        expect(input.password.value).toBe("");
        expect(input.submit.disabled).toBeTruthy();
    });

    test("empty input disables submit", () => {
        const { input } = setup();

        fireEvent.change(input.name, { target: { value: "username" } });
        expect(input.submit.disabled).toBeTruthy();

        fireEvent.change(input.name, { target: { value: "" } });
        fireEvent.change(input.password, { target: { value: "password123" } });
        expect(input.submit.disabled).toBeTruthy();
    });

    test("valid input enables submit", async () => {
        const { input } = setup();

        fireEvent.change(input.name, { target: { value: "username" } });
        fireEvent.change(input.password, { target: { value: "password123" } });

        expect(input.submit.disabled).toBeFalsy();
    });

    test("handles server error", async () => {
        const { input, onError, onSuccess } = setup();

        fireEvent.change(input.name, { target: { value: "error" } });
        fireEvent.change(input.password, { target: { value: "password123" } });
        fireEvent.click(input.submit);

        expect(input.submit.disabled).toBeTruthy();
        expect(input.submit.getAttribute("aria-busy")).toBe("true");
        await fast(() => expect(onError).toHaveBeenCalledWith(
            expect.objectContaining({
                message: "Server error",
                status: 500
            })
        ));
        expect(onSuccess).not.toHaveBeenCalled();
        expect(input.submit.disabled).toBeFalsy();
        expect(input.submit.getAttribute("aria-busy")).toBe("false");
        expect(input.name.value).toBe("error");
        expect(input.password.value).toBe("password123");
    });

    test("handles success", async () => {
        const { input, onError, onSuccess } = setup();

        fireEvent.change(input.name, { target: { value: "username" } });
        fireEvent.change(input.password, { target: { value: "password123" } });
        fireEvent.click(input.submit);

        expect(input.submit.disabled).toBeTruthy();
        expect(input.submit.getAttribute("aria-busy")).toBe("true");
        await fast(() => expect(onSuccess).toHaveBeenCalled());
        expect(onError).not.toHaveBeenCalled();
        expect(input.submit.disabled).toBeTruthy();
        expect(input.submit.getAttribute("aria-busy")).toBe("false");
    });
});