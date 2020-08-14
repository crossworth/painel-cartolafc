const changeLog = `
29 Julho
- Melhorado handling de errors de API
- Ajustado comportamento do menu selecionado
- Ajustado problema que causava resolver membros travar carregando
- Ajustado problema em "quotes" com UTF-8
- Adicionado função para lembrar a quantidade por páginas selecionado
- Adicionar informações de data de cache da página de membros
- Ajustado botões de ordenação da página de membros
- Adicionado suporte inicial de período na página de membros

30 Julho
- Adicionado suporte completo ao período na página de membros

1 Agosto
- Implementação inicial dos endpoints de tópicos
- Melhorado rotas da API
- Adicionado suporte a scroll na tabela para dispositivos com telas pequenas
- Implementado trava para evitar duas requests criarem o mesmo cache de uma rota de API lenta

2 Agosto
- Iniciado implementação de rotas de tópicos

3 Agosto
- Finalizado implementação de rotas de tópicos (porém performance ainda é baixa)

4 Agosto
- Melhorado performance rotas de tópicos e membros
- Finalizado implementação de rotas de tópicos (enquetes e perfils)

5 Agosto
- Adicionado suporte para retornar tópicos por data de atualização

6 Agosto
- Implementação de rotas de pesquisa

10 Agosto
- Implementação de rotas de ranking de tópicos
- Corrigido bug onde data de filtro não era aplicado de forma correta para 'última semana'

11 Agosto
- Adicionado suporte para 'último dia' como período.

14 Agosto
- Iniciados testes com FullTextSearch
- Ajustado problema onde cache podia armazenar erro em vez de dados
- Reduzido tempo de cache de 5horas para 1hora
- Ajustado referências ao ID da comunidade
- Adicionado listagem de ranking de tópicos
- Adicionado listagem de tópicos
- Adicionado suporte a 'short cache', para pesquisa
- Adicionado suporte inicial a pesquisa de tópicos
`

export {
  changeLog
}
