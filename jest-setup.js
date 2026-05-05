import React from 'react';
import './.config/jest-setup';
import { matchers } from './src/test/matchers';
import { MessageChannel } from 'worker_threads';

global.React = React;

// MessageChannel is used by @rc-component/select (and other libs) but is not
// in jsdom. Use Node's native implementation from worker_threads.
if (!global.MessageChannel) {
  global.MessageChannel = MessageChannel;
}

const mockIntersectionObserver = jest.fn().mockImplementation((callback) => ({
  observe: jest.fn().mockImplementation((elem) => {
    callback([{ target: elem, isIntersecting: true }]);
  }),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));
global.IntersectionObserver = mockIntersectionObserver;

// ResizeObserver is used by @rc-component/* (dropdowns, tooltips) but not in jsdom
global.ResizeObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));

expect.extend(matchers);

// jsdom doesn't apply user-agent stylesheets, so getComputedStyle(el).display returns ""
// for all elements. dom-accessibility-api uses this to decide whether to add a space
// separator between text nodes in accessible name computation. Without this patch, inline
// elements like <mark> and <span> are treated as block-level, splitting e.g. "value1-1"
// into "val  ue  1  -1" and breaking getByRole queries that match by accessible name.
const _origGetComputedStyle = window.getComputedStyle.bind(window);
const _inlineElements = new Set([
  'A', 'ABBR', 'ACRONYM', 'B', 'BDO', 'BIG', 'BR', 'BUTTON', 'CITE', 'CODE',
  'DFN', 'EM', 'I', 'IMG', 'INPUT', 'KBD', 'LABEL', 'MAP', 'MARK', 'OUTPUT',
  'Q', 'SAMP', 'SELECT', 'SMALL', 'SPAN', 'STRONG', 'S', 'SUB', 'SUP',
  'TEXTAREA', 'TIME', 'TT', 'U', 'VAR',
]);
Object.defineProperty(window, 'getComputedStyle', {
  value: (element, pseudo) => {
    const style = _origGetComputedStyle(element, pseudo);
    if (!pseudo && element && element.tagName && _inlineElements.has(element.tagName)) {
      const existingDisplay = element.style && element.style.display;
      if (!existingDisplay) {
        return new Proxy(style, {
          get(target, prop) {
            if (prop === 'display') {
              return 'inline';
            }
            if (prop === 'getPropertyValue') {
              return (name) => (name === 'display' ? 'inline' : target.getPropertyValue(name));
            }
            const val = target[prop];
            return typeof val === 'function' ? val.bind(target) : val;
          },
        });
      }
    }
    return style;
  },
  writable: true,
  configurable: true,
});
