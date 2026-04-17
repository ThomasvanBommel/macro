import { beforeAll, beforeEach, afterEach, afterAll, describe, test, expect, vi } from 'vitest'
import { http, HttpResponse, delay } from "msw";
import { setupServer } from "msw/node";

import { render, waitFor, fireEvent, cleanup } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

import RegisterForm from "./RegisterForm";

const server = setupServer(
    http.post("/api/register", async r => {
        const { name } = await r.request.json();

        if (name === "error")
            return HttpResponse.json({ error: "Server error" }, { status: 500 });

        // success, created
        return HttpResponse.json(null, { status: 201 });
    }),
);

function setup() {
    const onError = vi.fn();
    const onSuccess = vi.fn();

    const r = render(<RegisterForm onError={onError} onSuccess={onSuccess} />);

    const input = {
        name:     r.getByLabelText(/^username/i),
        password: r.getByLabelText(/^password/i),
        confirm:  r.getByLabelText(/^confirm/i),
        submit:   r.getByRole("button", { type: "submit" })
    };

    return { input, onError, onSuccess };
}

beforeAll(() => server.listen({ onUnhandledRequest: "error" }));
beforeEach(() => server.resetHandlers());
afterEach(() => cleanup());
afterAll(() => server.close());

describe("RegisterForm", () => {
    const fast = callback => waitFor(callback, { interval: 2 });

    test("initial state", () => {
        const { input } = setup();

        expect(input.name.value).toBe("");
        expect(input.password.value).toBe("");
        expect(input.confirm.value).toBe("");
        expect(input.submit.disabled).toBeTruthy();
        expect(input.password.getAttribute("aria-invalid")).toBeNull();
        expect(input.confirm.getAttribute("aria-invalid")).toBeNull();
    });

    test("passwords don't match", async () => {
        const { input } = setup();

        fireEvent.change(input.password, { target: { value: "password123" } });

        expect(input.password.getAttribute("aria-invalid")).toBe("true");
        expect(input.confirm.getAttribute("aria-invalid")).toBe("true");
        expect(input.submit.disabled).toBeTruthy();

        fireEvent.change(input.password, { target: { value: "" } });
        fireEvent.change(input.confirm, { target: { value: "password123" } });

        expect(input.password.getAttribute("aria-invalid")).toBe("true");
        expect(input.confirm.getAttribute("aria-invalid")).toBe("true");
        expect(input.submit.disabled).toBeTruthy();
    });

    test("aria-invalid should be null when both fields are empty", async () => {
        const { input } = setup();
        
        expect(input.password.getAttribute("aria-invalid")).toBeNull();
        expect(input.confirm.getAttribute("aria-invalid")).toBeNull();
        expect(input.submit.disabled).toBeTruthy();
    });

    test("username is required", async () => {
        const { input } = setup();
        
        fireEvent.change(input.password, { target: { value: "password123" } });
        fireEvent.change(input.confirm, { target: { value: "password123" } });

        expect(input.submit.disabled).toBeTruthy();
    });

    test("valid form enables submit", async () => {
        const { input } = setup();
        
        fireEvent.change(input.name, { target: { value: "newuser" } });
        fireEvent.change(input.password, { target: { value: "password123" } });
        fireEvent.change(input.confirm, { target: { value: "password123" } });

        expect(input.submit.disabled).toBeFalsy();
    });

    test("error response", async () => {
        const { input, onError, onSuccess } = setup();
        
        fireEvent.change(input.name, { target: { value: "error" } });
        fireEvent.change(input.password, { target: { value: "password123" } });
        fireEvent.change(input.confirm, { target: { value: "password123" } });
        fireEvent.click(input.submit);

        expect(input.submit.getAttribute("aria-busy")).toBe("true");
        await fast(() => expect(onError).toHaveBeenCalledWith(
            expect.objectContaining({
                message: "Server error",
                status: 500
            })
        ));
        expect(onSuccess).not.toHaveBeenCalled();
        expect(input.submit.getAttribute("aria-busy")).toBe("false");
    });

    test("successful registration", async () => {
        const { input, onError, onSuccess } = setup();
        
        fireEvent.change(input.name, { target: { value: "newuser" } });
        fireEvent.change(input.password, { target: { value: "password123" } });
        fireEvent.change(input.confirm, { target: { value: "password123" } });
        fireEvent.click(input.submit);

        expect(input.submit.getAttribute("aria-busy")).toBe("true");
        await fast(() => expect(onSuccess).toHaveBeenCalled());
        expect(onError).not.toHaveBeenCalled();
        expect(input.submit.getAttribute("aria-busy")).toBe("false");
    });
});