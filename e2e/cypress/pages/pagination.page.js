import { URL } from '../helpers';

const FIRST_PAGE_Btn = '[data-testid=firstPage]';
const PREV_PAGE_Btn = '[data-testid=prevPage]';
const NEXT_PAGE_Btn = '[data-testid=nextPage]';
const LAST_PAGE_Btn = '[data-testid=lastPage]';
const PAGINATION_STAGE = '[data-testid=paginationState]';

export default class PaginationPage {
    visitPage = (params) => {
        const url = URL.INDIVIDUALS + params;
        cy.log('navigating to %s', url);
        cy.visit(url);
        return this;
    };

    selectFirstPage = () => {
        cy.get(FIRST_PAGE_Btn).click();
        return this;
    };

    selectPrevPage = () => {
        cy.get(PREV_PAGE_Btn).click();
        return this;
    };

    selectNextPage = () => {
        cy.get(NEXT_PAGE_Btn).click();
        return this;
    };

    selectLastPage = () => {
        cy.get(LAST_PAGE_Btn).click();
        return this;
    };

    getPaginationState = () => {
        return cy.get(PAGINATION_STAGE);
    };
}
