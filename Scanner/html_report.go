package Scanner

import (
	"codeql_n1ght/Common"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// generateHTMLReport 生成HTML报告
func generateHTMLReport() error {
	// 读取SARIF文件
	sarifPath := "results.sarif"
	sarifData, err := os.ReadFile(sarifPath)
	if err != nil {
		Common.LogWarn("无法读取SARIF文件: %v", err)
		return nil // 不返回错误，因为这不是关键性失败
	}

	// 创建HTML报告
	htmlContent := generateHTMLTemplate(string(sarifData))

	// 写入HTML文件
	htmlPath := "scan_report.html"
	err = os.WriteFile(htmlPath, []byte(htmlContent), 0644)
	if err != nil {
		return fmt.Errorf("无法创建HTML报告: %v", err)
	}

	Common.LogInfo("HTML报告已生成: %s", htmlPath)

	// 自动打开HTML页面
	return openHTMLReport(htmlPath)
}

// getCurrentWorkingDir 获取当前工作目录
func getCurrentWorkingDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		Common.LogWarn("无法获取当前工作目录: %v", err)
		return ""
	}
	return pwd
}

// readFileContentsSafely 安全地读取文件内容
func readFileContentsSafely(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("无法读取文件: %v", err)
	}

	// 转义特殊字符以便安全嵌入到JavaScript中
	escaped := strings.ReplaceAll(string(content), "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "`", "\\`")
	escaped = strings.ReplaceAll(escaped, "$", "\\$")
	return escaped
}

// extractFilePathsFromSarif 从SARIF数据中提取所有文件路径
func extractFilePathsFromSarif(sarifData string) map[string]string {
	fileContents := make(map[string]string)

	// 这里简化处理，实际使用时会解析JSON
	// 但为了避免复杂的JSON解析，我们采用更直接的方式
	// 在JavaScript中处理时再动态读取文件

	return fileContents
}

// openHTMLReport 打开HTML报告
func openHTMLReport(htmlPath string) error {
	var cmd *exec.Cmd

	// 根据操作系统选择合适的命令
	switch {
	case strings.Contains(strings.ToLower(os.Getenv("OS")), "windows"):
		cmd = exec.Command("cmd", "/c", "start", htmlPath)
	case strings.Contains(strings.ToLower(os.Getenv("OSTYPE")), "darwin"):
		cmd = exec.Command("open", htmlPath)
	default:
		cmd = exec.Command("xdg-open", htmlPath)
	}

	if err := cmd.Start(); err != nil {
		Common.LogWarn("无法自动打开浏览器，请手动打开: %s", htmlPath)
		return nil
	}

	Common.LogInfo("已在浏览器中打开扫描报告: %s", htmlPath)
	return nil
}

// generateHTMLTemplate 生成HTML模板
func generateHTMLTemplate(sarifData string) string {
	// 确保sourceRootPath被正确初始化
	if sourceRootPath == "" {
		detectSourceRootPath()
	}

	// 获取当前工作目录并转换为正斜杠格式
	currentWorkingDir := strings.ReplaceAll(getCurrentWorkingDir(), "\\", "/")

	// 确保sourceRootPath使用正斜杠格式
	normalizedSource := strings.ReplaceAll(sourceRootPath, "\\", "/")

	// 确保数据库路径使用正斜杠格式
	normalizedDatabase := strings.ReplaceAll(Common.DatabasePath, "\\", "/")

	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CodeQL 扫描报告 - N1ght Scanner</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.3);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #2c3e50 0%, #3498db 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }
        
        .header p {
            font-size: 1.2em;
            opacity: 0.9;
        }
        
        .content {
            padding: 30px;
        }
        
        .section {
            margin-bottom: 30px;
            background: #f8f9fa;
            border-radius: 10px;
            padding: 20px;
            border-left: 5px solid #3498db;
        }
        
        .section h2 {
            color: #2c3e50;
            margin-bottom: 15px;
            font-size: 1.5em;
        }
        
        .sarif-content {
            background: #2c3e50;
            color: #ecf0f1;
            padding: 20px;
            border-radius: 8px;
            font-family: 'Courier New', monospace;
            overflow-x: auto;
            white-space: pre-wrap;
            line-height: 1.4;
            max-height: 600px;
            overflow-y: auto;
        }
        
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        
        .stat-card {
            background: white;
            border-radius: 10px;
            padding: 20px;
            text-align: center;
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
            border-top: 4px solid #3498db;
        }
        
        .stat-number {
            font-size: 2em;
            font-weight: bold;
            color: #3498db;
            margin-bottom: 5px;
        }
        
        .stat-label {
            color: #7f8c8d;
            font-size: 0.9em;
        }
        
        .footer {
            background: #34495e;
            color: white;
            padding: 20px;
            text-align: center;
        }
        
        .btn {
            display: inline-block;
            padding: 10px 20px;
            background: #3498db;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            margin: 5px;
            transition: background 0.3s;
        }
        
        .btn:hover {
            background: #2980b9;
        }
        
                 .timestamp {
             color: #95a5a6;
             font-size: 0.9em;
             margin-top: 10px;
         }
         
         .results-table {
             width: 100%;
             border-collapse: collapse;
             margin-top: 15px;
             background: white;
             border-radius: 8px;
             overflow: hidden;
             box-shadow: 0 4px 6px rgba(0,0,0,0.1);
             table-layout: fixed;
         }
         
         .results-table th {
             background: linear-gradient(135deg, #3498db, #2980b9);
             color: white;
             padding: 15px 10px;
             text-align: left;
             font-weight: 600;
             position: sticky;
             top: 0;
             z-index: 10;
         }
         
         .results-table th:nth-child(1) { width: 18%; } /* 查询规则 */
         .results-table th:nth-child(2) { width: 25%; } /* 文件路径 */
         .results-table th:nth-child(3) { width: 8%; }  /* 位置 */
         .results-table th:nth-child(4) { width: 10%; } /* 严重程度 */
         .results-table th:nth-child(5) { width: 20%; } /* 问题描述 */
         .results-table th:nth-child(6) { width: 19%; } /* 消息 */
         
         .results-table td {
             padding: 12px 10px;
             border-bottom: 1px solid #ecf0f1;
             vertical-align: top;
             word-wrap: break-word;
             word-break: break-word;
             overflow-wrap: break-word;
         }
         
         .results-table tbody tr:hover {
             background-color: #f8f9fa;
             transition: background-color 0.3s ease;
         }
         
         .results-table tbody tr:nth-child(even) {
             background-color: #fdfdfd;
         }
         
         .severity-high {
             background-color: #e74c3c !important;
             color: white;
             padding: 4px 8px;
             border-radius: 4px;
             font-weight: bold;
             text-align: center;
         }
         
         .severity-medium {
             background-color: #f39c12 !important;
             color: white;
             padding: 4px 8px;
             border-radius: 4px;
             font-weight: bold;
             text-align: center;
         }
         
         .severity-low {
             background-color: #27ae60 !important;
             color: white;
             padding: 4px 8px;
             border-radius: 4px;
             font-weight: bold;
             text-align: center;
         }
         
         .severity-info {
             background-color: #3498db !important;
             color: white;
             padding: 4px 8px;
             border-radius: 4px;
             font-weight: bold;
             text-align: center;
         }
         
         .file-path {
             font-family: 'Courier New', monospace;
             font-size: 0.85em;
             color: #2c3e50;
             background-color: #ecf0f1;
             padding: 2px 6px;
             border-radius: 3px;
             word-break: break-all;
             display: block;
             max-width: 100%;
             overflow-wrap: break-word;
             cursor: pointer;
             transition: all 0.3s ease;
         }
         
         .file-path:hover {
             background-color: #3498db;
             color: white;
             transform: scale(1.02);
         }
         
         .file-link {
             text-decoration: none;
             color: inherit;
         }
         
         .open-file-btn {
             background: #27ae60;
             color: white;
             border: none;
             padding: 2px 6px;
             border-radius: 3px;
             font-size: 0.7em;
             cursor: pointer;
             margin-left: 5px;
             transition: background 0.3s;
         }
         
         .open-file-btn:hover {
             background: #219a52;
         }
         
         .line-col {
             font-family: 'Courier New', monospace;
             font-weight: bold;
             color: #8e44ad;
             background-color: #f4f1ff;
             padding: 2px 6px;
             border-radius: 3px;
             text-align: center;
             white-space: nowrap;
         }
         
         .rule-id {
             font-family: 'Courier New', monospace;
             font-size: 0.9em;
             color: #27ae60;
             background-color: #eafaf1;
             padding: 2px 6px;
             border-radius: 3px;
         }
         
         .loading {
             text-align: center;
             font-style: italic;
             color: #7f8c8d;
             padding: 30px;
         }
         
         .no-results {
             text-align: center;
             color: #27ae60;
             padding: 30px;
             font-size: 1.1em;
         }
         
         .collapsible-section {
             margin-top: 15px;
         }
         
         .collapsible-btn {
             background: #3498db;
             color: white;
             border: none;
             padding: 10px 20px;
             border-radius: 5px;
             cursor: pointer;
             font-size: 0.9em;
             transition: background 0.3s;
         }
         
         .collapsible-btn:hover {
             background: #2980b9;
         }
         
         #results-table-container {
             max-height: 600px;
             overflow-y: auto;
             border-radius: 8px;
             border: 1px solid #ecf0f1;
         }
         
         /* 规则筛选器样式 */
         .filter-section {
             background: #f8f9fa;
             border-radius: 10px;
             padding: 20px;
             margin-bottom: 20px;
             border-left: 5px solid #27ae60;
         }
         
         .filter-controls {
             display: flex;
             flex-wrap: wrap;
             gap: 15px;
             align-items: center;
             margin-bottom: 15px;
         }
         
         .filter-group {
             display: flex;
             flex-direction: column;
             gap: 5px;
         }
         
         .filter-label {
             font-weight: 600;
             color: #2c3e50;
             font-size: 0.9em;
         }
         
         .filter-select {
             padding: 8px 12px;
             border: 2px solid #ecf0f1;
             border-radius: 5px;
             font-size: 0.9em;
             background: white;
             color: #2c3e50;
             cursor: pointer;
             transition: border-color 0.3s;
             min-width: 200px;
         }
         
         .filter-select:focus {
             outline: none;
             border-color: #3498db;
         }
         
         .filter-select:hover {
             border-color: #95a5a6;
         }
         
         .clear-filter-btn {
             background: #e74c3c;
             color: white;
             border: none;
             padding: 8px 16px;
             border-radius: 5px;
             font-size: 0.9em;
             cursor: pointer;
             transition: background 0.3s;
             height: fit-content;
             align-self: flex-end;
         }
         
         .clear-filter-btn:hover {
             background: #c0392b;
         }
         
         .filter-info {
             background: #e8f6f3;
             border: 1px solid #a2d9ce;
             color: #27ae60;
             padding: 10px 15px;
             border-radius: 5px;
             font-size: 0.9em;
         }
         
         .rule-option {
             padding: 5px 0;
         }
         
         .rule-option-content {
             display: flex;
             flex-direction: column;
             gap: 2px;
         }
         
         .rule-option-id {
             font-weight: bold;
             color: #27ae60;
         }
         
         .rule-option-name {
             font-size: 0.85em;
             color: #7f8c8d;
         }
        
        @media (max-width: 768px) {
            .stats {
                grid-template-columns: 1fr;
            }
            
            .header h1 {
                font-size: 2em;
            }
            
            .container {
                margin: 10px;
            }
            
            .filter-controls {
                flex-direction: column;
                align-items: stretch;
            }
            
            .filter-select {
                min-width: 100%;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 CodeQL 扫描报告</h1>
            <p>N1ght Scanner - 代码安全分析结果</p>
            <div class="timestamp">生成时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</div>
        </div>
        
        <div class="content">
                         <div class="stats">
                 <div class="stat-card">
                     <div class="stat-number" id="total-results">--</div>
                     <div class="stat-label">扫描结果总数</div>
                 </div>
                 <div class="stat-card">
                     <div class="stat-number" id="query-count">--</div>
                     <div class="stat-label">查询规则数量</div>
                 </div>
                 <div class="stat-card">
                     <div class="stat-number" id="file-count">--</div>
                     <div class="stat-label">涉及文件数量</div>
                 </div>
             </div>
            
            <div class="section">
                <h2>📊 扫描统计</h2>
                <p>数据库路径: <strong>` + Common.DatabasePath + `</strong></p>
                <p>查询库路径: <strong>` + Common.QLLibsPath + `</strong></p>
                <p>报告格式: <strong>SARIF v2.1.0</strong></p>
            </div>
             
            <div class="filter-section">
                <h2>🔍 结果筛选</h2>
                <div class="filter-controls">
                    <div class="filter-group">
                        <label class="filter-label" for="rule-filter">选择查询规则:</label>
                        <select id="rule-filter" class="filter-select" onchange="filterByRule()">
                            <option value="">所有规则</option>
                            <!-- 规则选项将由JavaScript动态生成 -->
                        </select>
                    </div>
                    <div class="filter-group">
                        <label class="filter-label" for="severity-filter">严重程度:</label>
                        <select id="severity-filter" class="filter-select" onchange="filterByRule()">
                            <option value="">所有级别</option>
                            <option value="error">🔴 高危</option>
                            <option value="warning">🟡 中危</option>
                            <option value="note">🟢 提示</option>
                            <option value="info">🔵 信息</option>
                        </select>
                    </div>
                    <button class="clear-filter-btn" onclick="clearAllFilters()">🔄 清除筛选</button>
                </div>
                <div id="filter-info" class="filter-info" style="display: none;">
                    正在显示筛选后的结果...
                </div>
            </div>
            
                         <div class="section">
                 <h2>📋 扫描结果详情</h2>
                 <div id="results-table-container">
                     <table id="results-table" class="results-table">
                         <thead>
                             <tr>
                                 <th>🔍 查询规则</th>
                                 <th>📁 文件路径</th>
                                 <th>📍 位置</th>
                                 <th>⚠️ 严重程度</th>
                                 <th>📝 问题描述</th>
                                 <th>💡 消息</th>
                             </tr>
                         </thead>
                         <tbody id="results-tbody">
                             <tr>
                                 <td colspan="6" class="loading">正在解析SARIF数据...</td>
                             </tr>
                         </tbody>
                     </table>
                 </div>
             </div>
             
             <div class="section">
                 <h2>📄 原始SARIF数据</h2>
                 <div class="collapsible-section">
                     <button class="collapsible-btn" onclick="toggleRawData()">显示/隐藏原始数据</button>
                     <div class="sarif-content" id="sarif-data" style="display: none;">` + sarifData + `</div>
                 </div>
             </div>
            
            <div class="section">
                <h2>🛠️ 操作选项</h2>
                <a href="#" class="btn" onclick="downloadSarif()">💾 下载 SARIF 文件</a>
                <a href="#" class="btn" onclick="printReport()">🖨️ 打印报告</a>
                <a href="#" class="btn" onclick="exportJson()">📤 导出 JSON</a>
            </div>
        </div>
        
        <div class="footer">
            <p>🚀 Powered by CodeQL N1ght Scanner | 🔒 专业代码安全分析工具</p>
            <p>如有问题请检查 SARIF 文件格式或联系技术支持</p>
        </div>
    </div>

         <script>
         let sarifData = null;
         let allResults = []; // 存储所有结果用于筛选
         let allRules = new Map(); // 存储所有规则信息
         
         // 解析SARIF数据并更新统计信息和表格
         function parseSarifData() {
             try {
                 const sarifText = document.getElementById('sarif-data').textContent;
                 sarifData = JSON.parse(sarifText);
                 
                 let totalResults = 0;
                 let queryCount = 0;
                 let fileSet = new Set();
                 allResults = []; // 重置结果数组
                 allRules.clear(); // 重置规则映射
                 
                 // 更新统计数据并收集所有结果
                 if (sarifData.runs && sarifData.runs.length > 0) {
                     sarifData.runs.forEach((run, runIndex) => {
                         // 收集规则信息
                         if (run.tool && run.tool.driver && run.tool.driver.rules) {
                             run.tool.driver.rules.forEach(rule => {
                                 allRules.set(rule.id, {
                                     id: rule.id,
                                     name: rule.name || rule.shortDescription?.text || rule.id,
                                     description: rule.shortDescription?.text || rule.fullDescription?.text || '无描述'
                                 });
                             });
                             queryCount += run.tool.driver.rules.length;
                         }
                         
                         if (run.results) {
                             totalResults += run.results.length;
                             
                             // 收集所有结果并添加运行信息
                             run.results.forEach((result, resultIndex) => {
                                 // 统计涉及的文件
                                 if (result.locations && result.locations.length > 0) {
                                     const location = result.locations[0];
                                     const physicalLocation = location.physicalLocation;
                                     if (physicalLocation && physicalLocation.artifactLocation) {
                                         fileSet.add(physicalLocation.artifactLocation.uri || '未知文件');
                                     }
                                 }
                                 
                                 // 添加到结果数组，包含运行信息
                                 allResults.push({
                                     result: result,
                                     run: run,
                                     runIndex: runIndex,
                                     resultIndex: resultIndex
                                 });
                             });
                         }
                     });
                 }
                 
                 document.getElementById('total-results').textContent = totalResults;
                 document.getElementById('query-count').textContent = queryCount;
                 document.getElementById('file-count').textContent = fileSet.size;
                 
                 // 生成规则筛选器选项
                 populateRuleFilter();
                 
                 // 生成结果表格
                 generateResultsTable();
                 
             } catch (e) {
                 console.log('SARIF 数据解析失败:', e);
                 document.getElementById('total-results').textContent = 'N/A';
                 document.getElementById('query-count').textContent = 'N/A';
                 document.getElementById('file-count').textContent = 'N/A';
                 
                 // 显示解析错误
                 const tbody = document.getElementById('results-tbody');
                 tbody.innerHTML = '<tr><td colspan="6" class="loading">SARIF 数据解析失败，请检查数据格式</td></tr>';
             }
         }
         
         // 填充规则筛选器选项
         function populateRuleFilter() {
             const ruleSelect = document.getElementById('rule-filter');
             
             // 清除现有选项（保留"所有规则"选项）
             while (ruleSelect.children.length > 1) {
                 ruleSelect.removeChild(ruleSelect.lastChild);
             }
             
             // 按规则ID排序
             const sortedRules = Array.from(allRules.values()).sort((a, b) => a.id.localeCompare(b.id));
             
             // 添加规则选项
             sortedRules.forEach(rule => {
                 const option = document.createElement('option');
                 option.value = rule.id;
                 option.innerHTML = ` + "`" + `${rule.id} - ${rule.name}` + "`" + `;
                 ruleSelect.appendChild(option);
             });
         }
         
         // 生成结果表格（支持筛选）
         function generateResultsTable(filteredResults = null) {
             const tbody = document.getElementById('results-tbody');
             tbody.innerHTML = '';
             
             // 使用筛选后的结果或所有结果
             const resultsToShow = filteredResults || allResults;
             
             if (resultsToShow.length === 0) {
                 if (allResults.length === 0) {
                     tbody.innerHTML = '<tr><td colspan="6" class="no-results">🎉 未发现安全问题！</td></tr>';
                 } else {
                     tbody.innerHTML = '<tr><td colspan="6" class="no-results">🔍 没有符合筛选条件的结果</td></tr>';
                 }
                 return;
             }
             
             resultsToShow.forEach(({ result, run }) => {
                 const row = document.createElement('tr');
                 
                 // 获取规则信息
                 const ruleId = result.ruleId || 'unknown';
                 const ruleName = getRuleName(run, ruleId);
                 
                 // 获取位置信息
                 const location = getLocationInfo(result);
                 
                 // 获取严重程度
                 const severity = result.level || 'note';
                 const severityClass = getSeverityClass(severity);
                 const severityText = getSeverityText(severity);
                 
                 // 获取消息
                 const message = result.message?.text || result.message?.markdown || '无描述';
                 
                 // 获取问题描述
                 const description = getRuleDescription(run, ruleId);
                 
                 row.innerHTML = ` + "`" + `
                     <td><span class="rule-id">${ruleId}</span><br/><small>${ruleName}</small></td>
                     <td>
                         <span class="file-path" onclick="openFileAtLocation('${location.file}', ${location.line}, ${location.column})" title="点击打开文件">${location.file}</span>
                         <button class="open-file-btn" onclick="openFileAtLocation('${location.file}', ${location.line}, ${location.column})" title="在编辑器中打开">📂</button>
                     </td>
                     <td><span class="line-col">${location.position}</span></td>
                     <td><span class="${severityClass}">${severityText}</span></td>
                     <td>${description}</td>
                     <td>${message}</td>
                 ` + "`" + `;
                 
                 tbody.appendChild(row);
             });
             
             // 更新筛选信息显示
             updateFilterInfo(resultsToShow.length);
         }
         
         // 按规则和严重程度筛选
         function filterByRule() {
             const ruleFilter = document.getElementById('rule-filter').value;
             const severityFilter = document.getElementById('severity-filter').value;
             
             let filteredResults = allResults;
             
             // 按规则筛选
             if (ruleFilter) {
                 filteredResults = filteredResults.filter(({ result }) => 
                     (result.ruleId || 'unknown') === ruleFilter
                 );
             }
             
             // 按严重程度筛选
             if (severityFilter) {
                 filteredResults = filteredResults.filter(({ result }) => 
                     (result.level || 'note') === severityFilter
                 );
             }
             
             // 更新统计信息
             updateFilteredStats(filteredResults);
             
             // 重新生成表格
             generateResultsTable(filteredResults);
         }
         
         // 清除所有筛选
         function clearAllFilters() {
             document.getElementById('rule-filter').value = '';
             document.getElementById('severity-filter').value = '';
             
             // 恢复原始统计信息
             updateOriginalStats();
             
             // 重新生成完整表格
             generateResultsTable();
             
             // 隐藏筛选信息
             document.getElementById('filter-info').style.display = 'none';
         }
         
         // 更新筛选后的统计信息
         function updateFilteredStats(filteredResults) {
             const totalResults = filteredResults.length;
             const fileSet = new Set();
             const ruleSet = new Set();
             
             filteredResults.forEach(({ result }) => {
                 // 统计涉及的文件
                 if (result.locations && result.locations.length > 0) {
                     const location = result.locations[0];
                     const physicalLocation = location.physicalLocation;
                     if (physicalLocation && physicalLocation.artifactLocation) {
                         fileSet.add(physicalLocation.artifactLocation.uri || '未知文件');
                     }
                 }
                 
                 // 统计涉及的规则
                 ruleSet.add(result.ruleId || 'unknown');
             });
             
             document.getElementById('total-results').textContent = totalResults;
             document.getElementById('query-count').textContent = ruleSet.size;
             document.getElementById('file-count').textContent = fileSet.size;
         }
         
         // 恢复原始统计信息
         function updateOriginalStats() {
             const totalResults = allResults.length;
             const fileSet = new Set();
             
             allResults.forEach(({ result }) => {
                 if (result.locations && result.locations.length > 0) {
                     const location = result.locations[0];
                     const physicalLocation = location.physicalLocation;
                     if (physicalLocation && physicalLocation.artifactLocation) {
                         fileSet.add(physicalLocation.artifactLocation.uri || '未知文件');
                     }
                 }
             });
             
             document.getElementById('total-results').textContent = totalResults;
             document.getElementById('query-count').textContent = allRules.size;
             document.getElementById('file-count').textContent = fileSet.size;
         }
         
         // 更新筛选信息显示
         function updateFilterInfo(filteredCount) {
             const filterInfo = document.getElementById('filter-info');
             const ruleFilter = document.getElementById('rule-filter').value;
             const severityFilter = document.getElementById('severity-filter').value;
             
             if (ruleFilter || severityFilter) {
                 const totalCount = allResults.length;
                 let filterText = ` + "`" + `正在显示 ${filteredCount} / ${totalCount} 条结果` + "`" + `;
                 
                 if (ruleFilter) {
                     const ruleName = allRules.get(ruleFilter)?.name || ruleFilter;
                     filterText += ` + "`" + ` (规则: ${ruleName})` + "`" + `;
                 }
                 
                 if (severityFilter) {
                     const severityText = getSeverityText(severityFilter);
                     filterText += ` + "`" + ` (级别: ${severityText})` + "`" + `;
                 }
                 
                 filterInfo.textContent = filterText;
                 filterInfo.style.display = 'block';
             } else {
                 filterInfo.style.display = 'none';
             }
         }
         
         // 获取规则名称
         function getRuleName(run, ruleId) {
             if (run.tool && run.tool.driver && run.tool.driver.rules) {
                 const rule = run.tool.driver.rules.find(r => r.id === ruleId);
                 return rule ? (rule.name || rule.shortDescription?.text || ruleId) : ruleId;
             }
             return ruleId;
         }
         
         // 获取规则描述
         function getRuleDescription(run, ruleId) {
             if (run.tool && run.tool.driver && run.tool.driver.rules) {
                 const rule = run.tool.driver.rules.find(r => r.id === ruleId);
                 return rule ? (rule.shortDescription?.text || rule.fullDescription?.text || '无描述') : '无描述';
             }
             return '无描述';
         }
         
         // 获取位置信息
         function getLocationInfo(result) {
             if (result.locations && result.locations.length > 0) {
                 const location = result.locations[0];
                 const physicalLocation = location.physicalLocation;
                 
                 if (physicalLocation) {
                     const file = physicalLocation.artifactLocation?.uri || '未知文件';
                     const region = physicalLocation.region;
                     
                     if (region) {
                         const startLine = region.startLine || 1;
                         const startColumn = region.startColumn || 1;
                         return {
                             file: file,
                             position: ` + "`" + `${startLine}:${startColumn}` + "`" + `,
                             line: startLine,
                             column: startColumn
                         };
                     }
                     
                     return {
                         file: file,
                         position: '未知位置',
                         line: 1,
                         column: 1
                     };
                 }
             }
             
             return {
                 file: '未知文件',
                 position: '未知位置',
                 line: 1,
                 column: 1
             };
         }
         
         // 获取严重程度样式类
         function getSeverityClass(severity) {
             switch(severity) {
                 case 'error': return 'severity-high';
                 case 'warning': return 'severity-medium';
                 case 'note': 
                 case 'info': return 'severity-low';
                 default: return 'severity-info';
             }
         }
         
         // 获取严重程度文本
         function getSeverityText(severity) {
             switch(severity) {
                 case 'error': return '🔴 高危';
                 case 'warning': return '🟡 中危';
                 case 'note': return '🟢 提示';
                 case 'info': return '🔵 信息';
                 default: return '⚪ 未知';
             }
         }
         
         // 切换原始数据显示
         function toggleRawData() {
             const rawData = document.getElementById('sarif-data');
             if (rawData.style.display === 'none') {
                 rawData.style.display = 'block';
             } else {
                 rawData.style.display = 'none';
             }
         }
         
         // 打开文件并定位到指定位置
         function openFileAtLocation(filePath, line, column) {
             // 在函数开始时声明变量，避免在catch块中出现未定义的问题
             let cleanPath = filePath;
             let actualLine = line;
             
             try {
                 // 如果路径包含行号列号信息（格式：path:line:column），则提取它们
                 const colonMatch = filePath.match(/^(.+?):(\d+):(\d+)$/);
                 if (colonMatch) {
                     cleanPath = colonMatch[1];      // 纯文件路径
                     actualLine = parseInt(colonMatch[2], 10);  // 行号
                 }
                 
                 // 清理路径中的冗余部分（如 /./ ）
                 cleanPath = cleanPath.replace(/\/\.\//g, '/');
                 
                 // 转换为绝对路径
                 const absolutePath = getAbsolutePath(cleanPath);
                 console.log('文件路径 (相对):', cleanPath);
                 console.log('文件路径 (绝对):', absolutePath);
                 
                 // 尝试多种编辑器打开方式
                 const editors = [
                     // VS Code
                     {
                         name: 'VS Code',
                         protocol: ` + "`" + `vscode://file/${absolutePath}:${actualLine}:${column}` + "`" + `,
                         command: ` + "`" + `code -g "${absolutePath}:${actualLine}:${column}"` + "`" + `
                     },
                     // Cursor 编辑器
                     {
                         name: 'Cursor',
                         protocol: ` + "`" + `cursor://file/${absolutePath}:${actualLine}:${column}` + "`" + `,
                         command: ` + "`" + `cursor -g "${absolutePath}:${actualLine}:${column}"` + "`" + `
                     },
                     // Sublime Text
                     {
                         name: 'Sublime Text',
                         command: ` + "`" + `subl "${absolutePath}:${actualLine}:${column}"` + "`" + `
                     },
                     // Notepad++
                     {
                         name: 'Notepad++',
                         command: ` + "`" + `notepad++ -n${actualLine} -c${column} "${absolutePath}"` + "`" + `
                     }
                 ];
                 
                 // 显示选择对话框
                 showEditorSelectionDialog(absolutePath, actualLine, column, editors);
                 
             } catch (error) {
                 console.error('打开文件失败:', error);
                 
                 // 备用方案：尝试直接打开文件（不定位）
                 try {
                     // 清理路径中的冗余部分（如 /./ ）
                     cleanPath = cleanPath.replace(/\/\.\//g, '/');
                     const absolutePath = getAbsolutePath(cleanPath);
                     const fileUrl = ` + "`" + `file:///${absolutePath}` + "`" + `;
                     window.open(fileUrl, '_blank');
                 } catch (fallbackError) {
                     console.error('备用方案也失败了:', fallbackError);
                     alert('无法打开文件: ' + filePath + '\\n错误: ' + error.message);
                 }
             }
         }
         
         // 获取绝对路径
         function getAbsolutePath(filePath) {
             // 如果已经是绝对路径，直接返回
             if (filePath.match(/^[a-zA-Z]:/)) {
                 return filePath.replace(/\\/g, '/');
             }
             
             // 构建基于动态检测的源码根目录的绝对路径
             // SARIF中的路径格式: src1/xxx.java
             const dbPath = '` + normalizedDatabase + `';
             const sourceRoot = '` + normalizedSource + `'; // 动态检测的源码根目录
             const currentWorkingDir = '` + currentWorkingDir + `'; // 当前工作目录
             
             let relativePath;
             if (sourceRoot) {
                 relativePath = dbPath + '/src/' + sourceRoot + '/' + filePath;
             } else {
                 // 后备方案：直接拼接
                 relativePath = dbPath + '/src/' + filePath;
             }
             
             // 转换为绝对路径
             let absolutePath = currentWorkingDir + '/' + relativePath;
             // 清理路径：去除冗余的斜杠和./部分
             absolutePath = absolutePath.replace(/\\/g, '/').replace(/\/+/g, '/').replace(/\/\.\//g, '/');
             // 处理开头的./
             if (absolutePath.includes('/.')) {
                 absolutePath = absolutePath.replace(/\/\./g, '');
             }
             return absolutePath;
         }

         // 显示编辑器选择对话框
         function showEditorSelectionDialog(filePath, line, column, editors) {
             const modal = document.createElement('div');
             modal.style.cssText = ` + "`" + `
                 position: fixed;
                 top: 0;
                 left: 0;
                 width: 100%;
                 height: 100%;
                 background: rgba(0,0,0,0.7);
                 display: flex;
                 justify-content: center;
                 align-items: center;
                 z-index: 1000;
             ` + "`" + `;
             
             const dialog = document.createElement('div');
             dialog.style.cssText = ` + "`" + `
                 background: white;
                 padding: 30px;
                 border-radius: 10px;
                 box-shadow: 0 10px 30px rgba(0,0,0,0.3);
                 max-width: 600px;
                 width: 90%;
             ` + "`" + `;
             
             dialog.innerHTML = ` + "`" + `
                 <h3 style="margin-top: 0; color: #2c3e50;">🚀 选择操作方式</h3>
                 <div style="background: #f8f9fa; padding: 15px; border-radius: 5px; margin-bottom: 20px;">
                     <p style="margin: 5px 0;"><strong>📄 文件:</strong> <code style="background: #e9ecef; padding: 2px 6px; border-radius: 3px;">${filePath}</code></p>
                     <p style="margin: 5px 0;"><strong>📍 位置:</strong> <code style="background: #e9ecef; padding: 2px 6px; border-radius: 3px;">第 ${line} 行，第 ${column} 列</code></p>
                 </div>
                 
                 <div style="margin-bottom: 20px;">
                     <h4 style="color: #34495e; margin-bottom: 10px;">📝 编辑器选择:</h4>
                     <div id="editor-buttons" style="margin-bottom: 15px;"></div>
                 </div>
                 
                 <div style="margin-bottom: 20px;">
                     <h4 style="color: #34495e; margin-bottom: 10px;">⚠️ 浏览器安全限制:</h4>
                     <div style="background: #fff3cd; border: 1px solid #ffeaa7; color: #856404; padding: 12px; border-radius: 5px; margin: 5px 0;">
                         由于浏览器的CORS安全策略，无法直接在网页中显示本地文件内容。<br/>
                         请使用下方的"📄 仅打开文件"按钮或"📋 复制路径"后在编辑器中打开。
                     </div>
                 </div>
                 
                 <div style="text-align: center;">
                     <button onclick="copyToClipboard('${filePath}:${line}:${column}')" 
                             style="background: #9b59b6; color: white; border: none; padding: 8px 16px; border-radius: 5px; margin: 5px; cursor: pointer;">
                         📋 复制路径
                     </button>
                     <button onclick="openFileOnly('${filePath}')" 
                             style="background: #34495e; color: white; border: none; padding: 8px 16px; border-radius: 5px; margin: 5px; cursor: pointer;">
                         📄 仅打开文件
                     </button>
                     <button onclick="closeModal()" 
                             style="background: #e74c3c; color: white; border: none; padding: 8px 16px; border-radius: 5px; margin: 5px; cursor: pointer;">
                         ❌ 取消
                     </button>
                 </div>
             ` + "`" + `;
             
             const buttonContainer = dialog.querySelector('#editor-buttons');
             
             editors.forEach(editor => {
                 const button = document.createElement('button');
                 button.style.cssText = ` + "`" + `
                     background: #3498db;
                     color: white;
                     border: none;
                     padding: 10px 15px;
                     border-radius: 5px;
                     margin: 5px;
                     cursor: pointer;
                     display: block;
                     width: 100%;
                     transition: background 0.3s;
                 ` + "`" + `;
                 button.textContent = ` + "`" + `📝 在 ${editor.name} 中打开` + "`" + `;
                 button.onmouseover = () => button.style.background = '#2980b9';
                 button.onmouseout = () => button.style.background = '#3498db';
                 button.onclick = () => {
                     tryOpenInEditor(editor, filePath, line, column);
                     closeModal();
                 };
                 buttonContainer.appendChild(button);
             });
             
             modal.appendChild(dialog);
             document.body.appendChild(modal);
             
             // 点击背景关闭
             modal.onclick = (e) => {
                 if (e.target === modal) closeModal();
             };
             
             window.closeModal = () => {
                 document.body.removeChild(modal);
             };
         }
         
         // 在弹窗中显示文件内容
         function showFileContent(filePath, targetLine) {
             // 分离文件路径和行号列号信息
             let cleanFilePath = filePath;
             let actualTargetLine = targetLine;
             
             // 如果路径包含行号列号信息（格式：path:line:column），则提取它们
             const colonMatch = filePath.match(/^(.+?):(\d+):(\d+)$/);
             if (colonMatch) {
                 cleanFilePath = colonMatch[1];
                 actualTargetLine = parseInt(colonMatch[2], 10);
             }
             
             // 清理路径中的冗余部分（如 /./ ）
             cleanFilePath = cleanFilePath.replace(/\/\.\//g, '/');
             
             console.log('清理后的文件路径:', cleanFilePath);
             console.log('目标行号:', actualTargetLine);
             
             // 显示文件信息和操作选项的弹窗
             showFileInfoModal(cleanFilePath, actualTargetLine);
         }
         
         // 显示文件信息和操作选项的弹窗
         function showFileInfoModal(filePath, targetLine) {
             const modal = document.createElement('div');
             modal.style.cssText = ` + "`" + `
                 position: fixed;
                 top: 0;
                 left: 0;
                 width: 100%;
                 height: 100%;
                 background: rgba(0,0,0,0.7);
                 display: flex;
                 justify-content: center;
                 align-items: center;
                 z-index: 1000;
             ` + "`" + `;
             
             const dialog = document.createElement('div');
             dialog.style.cssText = ` + "`" + `
                 background: white;
                 padding: 30px;
                 border-radius: 10px;
                 box-shadow: 0 10px 30px rgba(0,0,0,0.3);
                 max-width: 600px;
                 width: 90%;
             ` + "`" + `;
             
             dialog.innerHTML = ` + "`" + `
                 <h3 style="margin-top: 0; color: #2c3e50;">🚀 选择操作方式</h3>
                 <div style="background: #f8f9fa; padding: 15px; border-radius: 5px; margin-bottom: 20px;">
                     <p style="margin: 5px 0;"><strong>📄 文件:</strong> <code style="background: #e9ecef; padding: 2px 6px; border-radius: 3px;">${filePath}</code></p>
                     <p style="margin: 5px 0;"><strong>📍 位置:</strong> <code style="background: #e9ecef; padding: 2px 6px; border-radius: 3px;">第 ${targetLine} 行</code></p>
                 </div>
                 
                 <div style="margin-bottom: 20px;">
                     <h4 style="color: #34495e; margin-bottom: 10px;">⚠️ 浏览器安全限制:</h4>
                     <div style="background: #fff3cd; border: 1px solid #ffeaa7; color: #856404; padding: 12px; border-radius: 5px; margin: 5px 0;">
                         由于浏览器的CORS安全策略，无法直接在网页中显示本地文件内容。<br/>
                         请使用下方的"📄 仅打开文件"按钮或"📋 复制路径"后在编辑器中打开。
                     </div>
                 </div>
                 
                 <div style="text-align: center;">
                     <button onclick="copyToClipboard('${filePath}:${targetLine}')" 
                             style="background: #9b59b6; color: white; border: none; padding: 8px 16px; border-radius: 5px; margin: 5px; cursor: pointer;">
                         📋 复制路径
                     </button>
                     <button onclick="openFileOnly('${filePath}')" 
                             style="background: #34495e; color: white; border: none; padding: 8px 16px; border-radius: 5px; margin: 5px; cursor: pointer;">
                         📄 仅打开文件
                     </button>
                     <button onclick="closeModal()" 
                             style="background: #e74c3c; color: white; border: none; padding: 8px 16px; border-radius: 5px; margin: 5px; cursor: pointer;">
                         ❌ 取消
                     </button>
                 </div>
             ` + "`" + `;
             
             modal.appendChild(dialog);
             document.body.appendChild(modal);
             
             // 点击背景关闭
             modal.onclick = (e) => {
                 if (e.target === modal) closeModal();
             };
             
             window.closeModal = () => {
                 document.body.removeChild(modal);
             };
         }
         
         // 尝试在编辑器中打开
         function tryOpenInEditor(editor, filePath, line, column) {
             // 分离文件路径和行号列号信息（如果路径中包含的话）
             let cleanPath = filePath;
             let actualLine = line;
             let actualColumn = column;
             
             // 如果路径包含行号列号信息，则提取它们
             const colonMatch = filePath.match(/^(.+?):(\d+):(\d+)$/);
             if (colonMatch) {
                 cleanPath = colonMatch[1];
                 actualLine = parseInt(colonMatch[2], 10);
                 actualColumn = parseInt(colonMatch[3], 10);
             }
             
             // 清理路径中的冗余部分
             cleanPath = cleanPath.replace(/\/\.\//g, '/');
             
             // 转换为绝对路径
             const absolutePath = getAbsolutePath(cleanPath);
             
             // 首先尝试协议方式
             if (editor.protocol) {
                 // 使用绝对路径构建协议URL
                 const cleanProtocol = editor.protocol.replace(filePath, absolutePath).replace(':' + line + ':' + column, ':' + actualLine + ':' + actualColumn);
                 const link = document.createElement('a');
                 link.href = cleanProtocol;
                 link.click();
                 return;
             }
             
             // 显示命令行提示，使用绝对路径
             const cleanCommand = editor.command.replace(filePath, absolutePath).replace(':' + line + ':' + column, ':' + actualLine + ':' + actualColumn);
             showCommandInfo(cleanCommand);
         }
         
         // 显示命令行信息
         function showCommandInfo(command) {
             alert(` + "`" + `请在命令行中运行以下命令打开文件:\n\n${command}\n\n如果命令不存在，请确保对应编辑器已安装并添加到系统PATH中。` + "`" + `);
         }
         
         // 复制到剪贴板
         function copyToClipboard(text) {
             // 分离文件路径和行号列号信息
             let cleanPath = text;
             const colonMatch = text.match(/^(.+?):(\d+):(\d+)$/);
             if (colonMatch) {
                 cleanPath = colonMatch[1];
             }
             
             // 清理路径中的冗余部分
             cleanPath = cleanPath.replace(/\/\.\//g, '/');
             
             // 转换为绝对路径
             const absolutePath = getAbsolutePath(cleanPath);
             
             navigator.clipboard.writeText(absolutePath).then(() => {
                 alert('文件路径已复制到剪贴板！\n路径: ' + absolutePath);
             }).catch(() => {
                 // 备用方案
                 const textArea = document.createElement('textarea');
                 textArea.value = absolutePath;
                 document.body.appendChild(textArea);
                 textArea.select();
                 document.execCommand('copy');
                 document.body.removeChild(textArea);
                 alert('文件路径已复制到剪贴板！\n路径: ' + absolutePath);
             });
         }
         
         // 仅打开文件（不定位行号）
         function openFileOnly(filePath) {
             // 分离文件路径和行号列号信息
             let cleanPath = filePath;
             const colonMatch = filePath.match(/^(.+?):(\d+):(\d+)$/);
             if (colonMatch) {
                 cleanPath = colonMatch[1];
             }
             
             // 清理路径中的冗余部分
             cleanPath = cleanPath.replace(/\/\.\//g, '/');
             
             // 转换为绝对路径
             const absolutePath = getAbsolutePath(cleanPath);
             
             const fileUrl = ` + "`" + `file:///${absolutePath}` + "`" + `;
             window.open(fileUrl, '_blank');
         }
         
         function downloadSarif() {
             const sarifData = document.getElementById('sarif-data').textContent;
             const blob = new Blob([sarifData], { type: 'application/json' });
             const url = URL.createObjectURL(blob);
             const a = document.createElement('a');
             a.href = url;
             a.download = 'codeql_scan_results.sarif';
             document.body.appendChild(a);
             a.click();
             document.body.removeChild(a);
             URL.revokeObjectURL(url);
         }
         
         function printReport() {
             window.print();
         }
         
         function exportJson() {
             const sarifData = document.getElementById('sarif-data').textContent;
             const blob = new Blob([sarifData], { type: 'application/json' });
             const url = URL.createObjectURL(blob);
             const a = document.createElement('a');
             a.href = url;
             a.download = 'codeql_results.json';
             document.body.appendChild(a);
             a.click();
             document.body.removeChild(a);
             URL.revokeObjectURL(url);
         }
         
         // 页面加载时解析SARIF数据
         document.addEventListener('DOMContentLoaded', parseSarifData);
    </script>
</body>
</html>`
}
