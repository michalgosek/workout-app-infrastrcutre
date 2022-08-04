const classNames = {
    ACTIVE_LINK: 'nav-link router-link-exact-active',
    NON_ACTIVE_LINK: 'nav-link',
};

function GetLinkClassName(isActive: boolean): string {
    return isActive ? classNames.ACTIVE_LINK : classNames.NON_ACTIVE_LINK
}

export {
    GetLinkClassName
};