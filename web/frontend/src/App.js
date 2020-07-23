import React, { useState } from 'react'
import { BrowserRouter, Route, Switch } from 'react-router-dom'

import { AtlassianLogo, AtlassianIcon } from '@atlaskit/logo';
import Flag, { FlagGroup } from '@atlaskit/flag'

import HomePage from './pages/HomePage'
import ResolveProfileID from './pages/ResolveProfileID'
import { AtlassianNavigation, PrimaryButton } from '@atlaskit/atlassian-navigation'

const App = () => {
  const [flags, setFlags] = useState([])

  const onFlagDismissed = dismissedFlagId => {
    setFlags(flags.filter(flag => flag.id !== dismissedFlagId))
  }

  return (
    <div>
      <BrowserRouter>
        <Switch>
          <AtlassianNavigation
            label="site"
            primaryItems={[
              <PrimaryButton>Issues</PrimaryButton>,
            ]}
            logo={AtlassianLogo}
            renderProductHome={HomePage}
          >
            <Route path="/" component={HomePage}/>
            <Route path="/perfil" component={ResolveProfileID}/>
            <Route path="/perfil/:profileID" component={HomePage}/>
          </AtlassianNavigation>
        </Switch>
      </BrowserRouter>
      <FlagGroup onDismissed={onFlagDismissed}>
        {
          flags.map(flag => (
            <Flag
              id={flag.id}
              key={flag.id}
              title={flag.title}
              description={flag.description}
            />
          ))
        }
      </FlagGroup>
    </div>
  )
}

export default App
