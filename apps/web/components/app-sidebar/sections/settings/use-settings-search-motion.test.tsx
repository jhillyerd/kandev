import { cleanup, render } from "@testing-library/react";
import { afterEach, expect, it, vi } from "vitest";
import { useSettingsSearchMotion } from "./use-settings-search-motion";

type Handle = { cancel: () => void; effect: unknown; finish: () => void };

const originalAnimate = Object.getOwnPropertyDescriptor(Element.prototype, "animate");

afterEach(() => {
  cleanup();
  if (originalAnimate) Object.defineProperty(Element.prototype, "animate", originalAnimate);
  else Reflect.deleteProperty(Element.prototype, "animate");
});

function Results({ keys }: { keys: string[] }) {
  const ref = useSettingsSearchMotion(keys.join(","));
  return (
    <div ref={ref}>
      {keys.map((key) => (
        <div key={key} data-settings-search-motion-key={key}>
          {key}
        </div>
      ))}
    </div>
  );
}

it("releases an entrance animation's effect when it finishes", () => {
  const handles: Handle[] = [];
  const animate = vi.fn(() => {
    const listeners: (() => void)[] = [];
    const handle: Handle & { addEventListener: (type: string, fn: () => void) => void } = {
      cancel: vi.fn(),
      effect: {},
      addEventListener: (_type, fn) => listeners.push(fn),
      finish: () => listeners.forEach((fn) => fn()),
    };
    handles.push(handle);
    return handle;
  });
  Object.defineProperty(Element.prototype, "animate", { value: animate, configurable: true });
  const view = render(<Results keys={["a"]} />);
  view.rerender(<Results keys={["a", "b"]} />);
  const entrance = handles.at(-1);
  expect(entrance?.effect).not.toBeNull();
  entrance?.finish();
  expect(entrance?.effect).toBeNull();
});
