import React from 'react'
import { Link } from 'react-router-dom'
import Nav, { AkNavigationItem } from '@atlaskit/navigation'
import DashboardIcon from '@atlaskit/icon/glyph/dashboard'
import GearIcon from '@atlaskit/icon/glyph/settings'
import AtlassianIcon from '@atlaskit/icon/glyph/atlassian'
import ArrowleftIcon from '@atlaskit/icon/glyph/arrow-left'

const backIcon = <ArrowleftIcon label="Back icon" size="medium"/>
const globalPrimaryIcon = <AtlassianIcon label="Atlassian icon" size="xlarge"/>

const PageNavigation = () => {

  const navLinks = [
    ['/', 'Home', DashboardIcon],
    ['/settings', 'Settings', GearIcon],
  ]

  return (<Nav
    isOpen={true}
    width={300}
    onResize={() => {}}
    globalPrimaryIcon={globalPrimaryIcon}
    globalPrimaryItemHref="/">
    {
      navLinks.map(link => {
        const [url, title, Icon] = link
        return (
          <Link key={url} to={url}>
            <AkNavigationItem
              icon={<Icon label={title} size="medium"/>}
              text={title}
              isSelected={false}
            />
          </Link>
        )
      }, this)
    }
  </Nav>)
}

export default PageNavigation
