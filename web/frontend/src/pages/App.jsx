import React from 'react'
import { BrowserRouter as Router, Link, Route, Switch } from 'react-router-dom'

import { ConfigProvider, Layout, Menu } from 'antd'
import ptBR from 'antd/es/locale/pt_BR'

import HomeOutlined from '@ant-design/icons/lib/icons/HomeOutlined'
import TeamOutlined from '@ant-design/icons/lib/icons/TeamOutlined'
import LinkOutlined from '@ant-design/icons/lib/icons/LinkOutlined'
import SearchOutlined from '@ant-design/icons/lib/icons/SearchOutlined'
import LogoutOutlined from '@ant-design/icons/lib/icons/LogoutOutlined'
import SettingOutlined from '@ant-design/icons/lib/icons/SettingOutlined'
import OrderedListOutlined from '@ant-design/icons/lib/icons/OrderedListOutlined'
import UnorderedListOutlined from '@ant-design/icons/lib/icons/UnorderedListOutlined'

import Profile from './Profile'
import ProfileList from './ProfileList'
import PageNotFound from './PageNotFound'
import ProfileTopics from './ProfileTopics'
import ProfileResolve from './ProfileResolve'
import ProfileComments from './ProfileComments'
import ProfileNotFound from './ProfileNotFound'
import TopicList from './TopicList'
import TopicRankingList from './TopicRankingList'
import { changeLog } from '../changelog'
import TopicSearch from './TopicSearch'
import Settings from './Settings'

const { Sider, Content } = Layout

export default () => {
  const selectedMenu = window.location.pathname.substr(0)

  return (
    <ConfigProvider locale={ptBR}>
      <Router>
        <Layout style={{ minHeight: '100vh' }}>
          <Sider
            className="side-area"
            breakpoint="lg"
            width="240px">
            <div className="logo">
              <svg fill="none" width="180" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 445 112">
                <path
                  d="M60.904 51.8926c-1.178-.3778-2.346-.6165-3.4978-.7189-1.1519-.1024-2.6264-.0626-4.4246.1243l-14.7554 1.5343c-1.4615.1512-2.5581.5032-3.2868 1.0511-.7317.5499-1.2664 1.5094-1.6032 2.8787L20.7419 98.5636c-.1688.5171-.2543.9715-.2543 1.3593 0 1.3901 1.1238 1.9691 3.3722 1.7341L315.938 71.5574l-4.868 19.45L10.7077 111.868c-3.5964.372-6.2821-.043-8.0511-1.25C.8876 109.415 0 107.168 0 103.888c0-1.779.3377-4.0068 1.0122-6.6936l10.9097-34.7863c1.4031-5.7624 3.5401-10.0014 6.4067-12.718 2.8666-2.7176 6.6037-4.3156 11.2133-4.7929l32.3753-3.3679-1.0132 10.3633z"
                  fill="#1890ff"/>
                <path
                  d="M84.44 53.002c3.635-.38 6.305-.032 8.012 1.047 1.705 1.08 2.558 3.1 2.558 6.06 0 1.79-.336 4.03-1.007 6.73l-4.277 17.64c-.393 1.66-.588 3.08-.588 4.25 0 1.12.225 2.15.672 3.11l-19.797 2.08.084-4.37c-1.063 1.9-2.18 3.19-3.355 3.87-1.173.68-3.158 1.17-5.953 1.47l-14.1 1.48c-3.3.34-5.66.06-7.09-.86-1.43-.92-2.14-2.52-2.14-4.83 0-1.515.25-3.308.75-5.38l.75-2.94c1.06-4.43 2.63-7.667 4.69-9.708 2.06-2.041 5.28-3.29 9.64-3.75l18.03-1.894c1.4-.143 2.35-.41 2.852-.79.502-.38.85-.943 1.05-1.68.192-.734.29-1.184.29-1.35 0-1.045-.98-1.465-2.937-1.26L52.865 64c-2.96.31-5.73 1.08-8.3 2.3l3.776-9.54L84.44 53v.002zM63.974 75.54c-1.566.163-2.7.49-3.396.983-.7.492-1.245 1.46-1.637 2.898l-1.09 4.65c-.16.64-.25 1.15-.25 1.54 0 .73.25 1.19.76 1.39.51.2 1.46.23 2.85.08l4.62-.48c1.9-.2 3.16-.62 3.78-1.27.62-.65 1.23-2.29 1.85-4.93l1.35-5.76-8.81.93-.027-.031zm57.96-26.477l-.92 6.724c1.677-2.733 3.325-4.65 4.95-5.75 1.62-1.103 3.94-1.812 6.96-2.13l9.645-1.013-3.19 12.877c-1.12-.55-1.87-.76-2.98-.89-1.08-.12-2.87-.05-5.66.25-3.81.4-6.42 1.13-7.85 2.17-1.43 1.05-2.48 2.95-3.15 5.7l-3.19 13c-.56 2.36-.84 4.17-.84 5.46 0 .62.05 1.13.17 1.54.11.41.3 1.1.58 2.08l-24.15 2.53c1.17-1.52 2.05-2.84 2.64-3.97.58-1.12 1.18-2.97 1.8-5.55l5.28-21.86c.56-2.35.85-3.93.87-5.27.01-.61-.14-1.63-.25-2.04-.12-.41-.26-.78-.54-1.76l19.79-2.08.035-.018zm27.566-10.95l19.964-2.6-2.182 8.784 15.77-1.658-1.93 7.92-15.77 1.65-5.2 21.34c-.167.69-.25 1.23-.25 1.62 0 .61.236 1.01.71 1.18.474.17 1.386.19 2.728.05l5.672-.73c1.678-.18 3.088-.48 4.235-.91 1.146-.43 2.224-.89 3.23-1.39l-3.522 9.6-23.54 2.6c-2.85.3-5.114-.02-6.792-.96s-2.518-2.81-2.518-5.61c0-1.01.028-1.77.085-2.28.054-.507.25-1.37.585-2.58l4.782-19.874 1.93-7.923 2.014-8.26-.001.031zM233.968 64.8c-1.064 4.42-2.63 7.645-4.697 9.68-2.07 2.034-5.28 3.28-9.64 3.74l-30.86 3.243c-3.3.348-5.66.062-7.09-.85-1.42-.91-2.14-2.514-2.14-4.807 0-1.512.25-3.297.76-5.363l3.95-16.017c1.06-4.418 2.63-7.643 4.7-9.678 2.07-2.034 5.29-3.283 9.65-3.74l30.87-3.245c3.3-.347 5.67-.062 7.09.85 1.42.912 2.14 2.515 2.14 4.81 0 1.508-.25 3.295-.75 5.36L234 64.8h-.032zm-34.14 4.595c-.167.635-.25 1.145-.25 1.537 0 1.12.78 1.596 2.348 1.43l8.22-.863c2.068-.22 3.327-1.41 3.775-3.59l4.45-18.08c.17-.64.25-1.15.25-1.54 0-1.12-.84-1.59-2.52-1.41l-8.22.86c-1.95.2-3.16 1.39-3.6 3.56l-4.44 18.08-.013.016zm48.232-38.79c.557-2.35.877-4.093.877-5.38 0-.615-.096-1.2-.207-1.61-.113-.407-.308-1.1-.587-2.076l21.976-2.31-11.16 45.79c-.56 2.35-.84 4.17-.84 5.45 0 .61.05 1.12.17 1.53.11.41.3 1.1.58 2.08l-24.16 2.54c1.17-1.52 2.05-2.85 2.64-3.97.59-1.13 1.19-2.98 1.8-5.56l8.89-36.5.021.016zm63.662-1.49c3.635-.38 6.305-.033 8.012 1.046 1.703 1.08 2.56 3.11 2.56 6.07 0 1.79-.337 4.04-1.01 6.73L317.01 60.6c-.393 1.663-1.383 5.134-1.867 7.58l-17.845 1.874.086-4.37c-1.063 1.9-2.183 3.192-3.354 3.874-1.174.686-3.16 1.172-5.957 1.467l-14.093 1.48c-3.3.346-5.662.064-7.088-.855-1.427-.914-2.14-2.52-2.14-4.822 0-1.514.254-3.307.756-5.38l.756-2.94c1.06-4.432 2.625-7.666 4.695-9.708 2.06-2.04 5.28-3.288 9.64-3.75l18.03-1.894c1.39-.146 2.34-.41 2.85-.794.5-.384.85-.94 1.05-1.677.19-.736.29-1.186.29-1.352 0-1.045-.98-1.464-2.94-1.26l-19.71 2.072c-2.97.312-5.74 1.08-8.31 2.3l3.77-9.54 36.07-3.79h.023zM291.255 51.65c-1.566.165-2.7.492-3.396.985-.71.494-1.25 1.46-1.64 2.898l-1.09 4.644c-.17.633-.25 1.145-.25 1.535 0 .73.25 1.19.75 1.39.5.2 1.45.224 2.85.077l4.61-.49c1.9-.2 3.16-.63 3.77-1.28.61-.65 1.23-2.3 1.84-4.94l1.34-5.76-8.81.92.026.021zm50.905 25.47c-.78 3.136-1.17 5.645-1.17 7.53 0 1.04.163 2.72.487 3.855l-27.688 2.316c.75-3.6 19.99-78.81 19.92-79.13l53.39-5.54-3.54 11.68c-.62-.26-1.06-.36-2.7-.5-1.69-.14-3.61.12-5.95.37l-19.4 2.04-3.61 15 16.96-1.78c2.53-.26 4.46-.51 5.8-.75 1.33-.23 3.1-.6 5.31-1.1l-3.15 13.22c-.93-.37-1.33-.59-3.72-.87-1.65-.19-4.6-.13-6.87.11l-16.97 1.78-7.136 31.79.037-.021zm99.454-65.26c-1.193-.322-1.412-.338-2.553-.432-1.33-.11-3.07-.002-5.15.22l-17.06 1.817c-1.69.18-2.96.597-3.8 1.246-.84.65-1.46 1.79-1.85 3.41l-11.42 50.29c-.19.62-.29 1.15-.29 1.61 0 1.65 1.3 2.34 3.9 2.06l13.45-1.43c4.62-.49 8.28-1.23 9.87-1.69 2-.58 2.91-1.05 5.05-2.27L426.6 81l-38.26 3.25c-4.16.44-7.414-.114-9.46-1.544-2.05-1.427-3.073-4.087-3.073-7.976 0-2.11 3.428-19.94 4.21-23.12l6.435-26.784c1.625-6.83 4.096-11.853 7.41-15.073 3.314-3.22 7.637-5.11 12.967-5.68l37.72-3.697-2.92 11.5-.015-.016z"
                  fill="#1890ff"/>
              </svg>
            </div>
            <Menu
              mode="inline"
              defaultSelectedKeys={selectedMenu}>
              <Menu.Item key="">
                <Link to="/">
                  <HomeOutlined/>
                  <span>Home</span>
                </Link>
              </Menu.Item>

              <Menu.Item key="/topicos/todos">
                <Link to="/topicos/todos">
                  <UnorderedListOutlined/>
                  <span>Tópicos</span>
                </Link>
              </Menu.Item>

              <Menu.Item key="/topicos/ranking">
                <Link to="/topicos/ranking">
                  <OrderedListOutlined/>
                  <span>Ranking Tópicos</span>
                </Link>
              </Menu.Item>

              <Menu.Item key="/topicos/pesquisa">
                <Link to="/topicos/pesquisa">
                  <SearchOutlined/>
                  <span>Pesquisa Tópicos</span>
                </Link>
              </Menu.Item>

              <Menu.Item key="perfil/todos">
                <Link to="/perfil/todos">
                  <TeamOutlined/>
                  <span>Membros</span>
                </Link>
              </Menu.Item>

              <Menu.Item key="resolver">
                <Link to="/resolver">
                  <LinkOutlined/>
                  <span>Resolver nome/link</span>
                </Link>
              </Menu.Item>

              {
                window.User.type === 'super_admin' && <Menu.Item key="configuracoes">
                  <Link to="/configuracoes">
                    <LogoutOutlined/>
                    <span>Configurações</span>
                  </Link>
                </Menu.Item>
              }

              <Menu.Item key="logout">
                <a href="/logout">
                  <SettingOutlined/>
                  <span>Sair</span>
                </a>
              </Menu.Item>
            </Menu>
          </Sider>

          <Content className="main-content">
            <Switch>
              <Route path='/' exact>
                <h3>WIP: Preview!</h3>
                Faltando: <br/>
                - <s>Lista de tópicos</s><br/>
                - Busca por título de tópico, conteúdo comentário, data, membro, número de comentários<br/>
                - Reconstituir tópicos apagados (recriar uma visualização de um tópico apagado)<br/>
                - Opções de exportar dados para Excel, CSV<br/>
                - <s>Verificações se tópico foi apagado</s> Adicionar tabela de metadados de tópico com
                se-deletado/número-comentários<br/>
                - Verificações se membro foi/está banido<br/>
                - <s>Na lista de membros, adicionar filtros por períodos (membro com mais tópicos na semana, mês e
                geral)</s><br/>
                - Melhorar performance da rota de status de membros (atualmente leva ~6s) (adicionar tabela de
                metadados?)<br/>
                - <s>Melhorar performance da rota de tópicos (contagem de comentários leva muito tempo +1m)</s><br/>
                - Adicionar login com VK e remover BasicAuth<br/>
                - Separar conteúdo membro/administrador<br/>
                - Melhorar forma de contar conteúdo no banco e dados (atualmente é o que demora mais nas queries)<br/>
                - Adicionar soft-delete<br/>
                - <s>Adicionar página de filtros a tópicos (mais comentários geral, mês, semana)</s><br/>
                - Link de tópico aleatório<br/>

                <br/><br/><br/><br/><br/>
                <h3>ChangeLog</h3>
                <div style={{ whiteSpace: 'pre' }} dangerouslySetInnerHTML={{ __html: changeLog }}/>
              </Route>

              <Route path='/resolver/:name?' render={(props) => <ProfileResolve {...props}/>}>
              </Route>

              {
                window.User.type === 'super_admin' && <Route path='/configuracoes' exact>
                  <Settings/>
                </Route>
              }

              <Route path='/topicos/pesquisa' exact>
                <TopicSearch/>
              </Route>

              <Route path='/topicos/todos' exact>
                <TopicList/>
              </Route>

              <Route path='/topicos/ranking' exact>
                <TopicRankingList/>
              </Route>

              <Route path='/perfil/nao-encontrado' exact>
                <ProfileNotFound/>
              </Route>

              <Route path='/perfil/todos' exact>
                <ProfileList/>
              </Route>

              <Route path='/perfil/:id' exact render={(props) => <Profile {...props}/>}>
              </Route>

              <Route path='/perfil/:id/topicos' exact render={(props) => <ProfileTopics {...props}/>}>
              </Route>

              <Route path='/perfil/:id/comentarios' exact render={(props) => <ProfileComments {...props}/>}>
              </Route>

              <Route component={PageNotFound}/>
            </Switch>
          </Content>
        </Layout>
      </Router>
    </ConfigProvider>
  )
}
