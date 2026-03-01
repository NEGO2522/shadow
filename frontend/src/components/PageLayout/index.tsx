import clsx from 'clsx';
import { ReactNode } from 'react';
import { twMerge } from 'tailwind-merge';

/**
 * This component is a simple page layout component to help with design consistency
 * Feel free to modify this component to fit your needs
 */
export const Page = (props: { children: ReactNode; className?: string }) => {
  return (
    <div className={twMerge(clsx('app-screen', props.className))}>
      {props.children}
    </div>
  );
};

const Header = (props: { children: ReactNode; className?: string }) => {
  return (
    <header
      className={twMerge(
        'screen-header border-b border-white/5 bg-[#0f1011]/80 backdrop-blur-md sticky top-0 z-50',
        clsx(props.className),
      )}
    >
      {props.children}
    </header>
  );
};

const Main = (props: { children: ReactNode; className?: string }) => {
  return (
    <main
      className={twMerge(
        clsx('grow overflow-y-auto p-6 pt-3 no-scrollbar', props.className),
      )}
    >
      {props.children}
    </main>
  );
};

const Footer = (props: { children: ReactNode; className?: string }) => {
  return (
    <footer className={twMerge('bottom-nav-v2', clsx(props.className))}>
      {props.children}
    </footer>
  );
};

Page.Header = Header;
Page.Main = Main;
Page.Footer = Footer;
