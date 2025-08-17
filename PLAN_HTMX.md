# PLAN_HTMX.md

# HTMX Implementation Plan


# Overview

This plan outlines the incremental implementation of HTMX support for the snippetbox application. The approach follows TDD principles with small, testable steps that build upon each other.



# Kanban Board Structure



### 🟢 BACKLOG - Ready Tasks

- [ ] Setup basic HTMX integration in templates

- [ ] Implement HTMX-aware handlers

- [ ] Add HTMX JavaScript functionality

- [ ] Create test coverage for HTMX features

- [ ] Implement form submission with HTMX

- [ ] Add loading states and visual feedback

- [ ] Handle error responses with HTMX



### 🟡 IN PROGRESS - Active Development

- [ ] (Empty initially)



### 🟢 DONE - Completed Tasks

- [x] Setup basic HTMX integration in templates

  - [x] Verify HTMX library is loaded in base layout
  - [x] Add basic HTMX attributes to navigation elements
  - [x] Create test for HTMX-enabled template rendering
  - [x] Ensure existing functionality remains intact


# Detailed Implementation Plan



### Phase 1: Foundation Setup


#### Task 1: Setup basic HTMX integration in templates

**Sub-tasks:**

- [x] Verify HTMX library is loaded in base layout

  - Validation: Check that `https://unpkg.com/htmx.org@1.9.6` is included in the `<head>` section of `base.layout.tmpl`

- [x] Add basic HTMX attributes to navigation elements

  - Validation: Verify that navigation links in `base.layout.tmpl` have appropriate HTMX attributes (e.g., `hx-get`, `hx-target`)

- [x] Create test for HTMX-enabled template rendering

  - Validation: Write a test that confirms templates render correctly with HTMX attributes present

- [x] Ensure existing functionality remains intact

  - Validation: Run all existing tests to confirm no regression in non-HTMX functionality


#### Task 2: Implement HTMX-aware handlers

**Sub-tasks:**

- [ ] Modify handler functions to detect HTMX requests

  - Validation: Create a helper function that checks `r.Header.Get("HX-Request") == "true"` and returns boolean

- [ ] Add conditional response logic based on HTMX flag

  - Validation: Verify handlers return different content (full page vs partial) based on HTMX request detection

- [ ] Create helper function to check for HTMX requests

  - Validation: Test that the helper function correctly identifies HTMX requests and non-HTMX requests

- [ ] Update existing tests to account for new behavior

  - Validation: Run all existing handler tests to ensure they pass with new HTMX-aware logic


### Phase 2: Core Functionality


#### Task 3: Implement form submission with HTMX

**Sub-tasks:**

- [ ] Add HTMX attributes to create snippet form

  - Validation: Confirm `create.page.tmpl` has `hx-post` attribute pointing to `/snippet/create` and `hx-target` pointing to the form container

- [ ] Modify createSnippet handler to handle HTMX submissions

  - Validation: Test that when HTMX request is detected, handler returns appropriate partial HTML instead of full page redirect

- [ ] Return partial HTML on successful form submission

  - Validation: Verify successful form submission returns updated snippet list or success message in partial format

- [ ] Handle validation errors in HTMX context

  - Validation: Confirm validation errors are returned as partial HTML with error messages displayed appropriately


#### Task 4: Add loading states and visual feedback

**Sub-tasks:**

- [ ] Implement loading indicators for HTMX operations

  - Validation: Add CSS classes for loading states and ensure they're applied during HTMX requests

- [ ] Add CSS styling for loading states

  - Validation: Verify CSS rules exist for loading indicators in static stylesheets

- [ ] Create JavaScript to show/hide loading indicators

  - Validation: Test that JavaScript properly toggles loading state classes on HTMX events

- [ ] Test loading state behavior with HTMX

  - Validation: Run integration tests confirming loading states appear during HTMX operations


### Phase 3: Advanced Features

#### Task 5: Handle error responses with HTMX

**Sub-tasks:**

- [ ] Implement proper error handling in HTMX context

  - Validation: Verify error responses return appropriate partial HTML content instead of full error pages

- [ ] Return appropriate error messages in partial HTML

  - Validation: Test that form validation errors display correctly in partial HTML responses

- [ ] Add visual feedback for form errors

  - Validation: Confirm error styling is applied to form elements in HTMX context

- [ ] Test error scenarios with HTMX

  - Validation: Run tests covering various error conditions with HTMX requests


#### Task 6: Create comprehensive test coverage

**Sub-tasks:**

- [ ] Add integration tests for HTMX functionality

  - Validation: Write tests that simulate full HTMX workflows from request to response

- [ ] Test both regular and HTMX request paths

  - Validation: Ensure all routes work correctly for both standard and HTMX requests

- [ ] Verify partial HTML responses work correctly

  - Validation: Confirm partial HTML fragments are properly returned for HTMX requests

- [ ] Ensure existing tests still pass

  - Validation: Run complete test suite to verify no regressions


## Implementation Details



### Template Modifications Required:

1. **base.layout.tmpl**: Already includes HTMX library

2. **create.page.tmpl**: Add `hx-post` to form with appropriate target

3. **home.page.tmpl**: Add HTMX attributes for dynamic content updates


### Handler Modifications Required:

1. **createSnippet handler**: Check for HTMX requests and return partial HTML

2. **home handler**: Support returning partial content when requested via HTMX

3. **Helper functions**: Add HTMX detection logic


### JavaScript Enhancements:

1. **main.js**: Add loading state management

2. **Event handling**: Handle HTMX events appropriately


## Testing Strategy



### Unit Tests:

- [ ] Test HTMX request detection

  - Validation: Create unit test that calls the HTMX detection helper with different request headers

- [ ] Test partial HTML response generation

  - Validation: Verify handlers return correct content types and structures for HTMX requests

- [ ] Test form submission with HTMX context

  - Validation: Test form handling logic specifically for HTMX scenarios


### Integration Tests:

- [ ] Test full HTMX workflow from form submission to response

  - Validation: Simulate complete HTMX interaction including request, processing, and response

- [ ] Verify loading states work correctly

  - Validation: Confirm visual feedback appears during HTMX operations

- [ ] Test error handling in HTMX context

  - Validation: Validate error responses are properly formatted for HTMX requests


## Success Criteria



### Minimum Viable Product:

1. Form submissions work via HTMX

2. Partial HTML responses are returned for HTMX requests

3. Basic loading indicators are implemented

4. Existing functionality remains intact


### Enhanced Features:

1. Proper error handling with HTMX

2. Comprehensive test coverage

3. Smooth user experience with visual feedback


## Dependencies

- HTMX library already included in base layout

- Existing form validation logic (forms package)

- Current template rendering system

- Session management for authentication


## Risk Mitigation

- Maintain backward compatibility with existing non-HTMX requests

- Ensure all existing tests continue to pass

- Implement incremental changes with frequent testing

- Keep template modifications minimal and focused


## Next Steps

1. Start with Phase 1 tasks (Foundation Setup)

2. Implement basic HTMX integration in templates

3. Modify handlers to detect and respond to HTMX requests

4. Test the basic functionality

5. Proceed to Phase 2 and 3 tasks incrementally